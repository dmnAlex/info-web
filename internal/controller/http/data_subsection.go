package http

import (
	"encoding/csv"
	"log"
	"net/http"

	"github.com/dsnikitin/info-web/internal/entity"
	"github.com/dsnikitin/info-web/internal/pkg/tools"
	"github.com/dsnikitin/info-web/internal/template"
	"github.com/gin-gonic/gin"
)

type DataSubsectionUseCase[E entity.Entity] interface {
	GetAll() ([]E, error)
	Create(e *E) error
	Update(e *E) error
	Delete(id string) error
	Export() ([][]string, error)
	// Import()
}

type DataSubsection[E entity.Entity] struct {
	uc DataSubsectionUseCase[E]
}

func NewDataSubsection[E entity.Entity](uc DataSubsectionUseCase[E]) *DataSubsection[E] {
	return &DataSubsection[E]{uc: uc}
}

func (ds *DataSubsection[E]) GetAllEntities(ctx *gin.Context) {
	if entities, err := ds.uc.GetAll(); err != nil {
		ctx.HTML(http.StatusInternalServerError, template.InternalServerError, "")
	} else {
		content := ds.preparePageContent(entities)
		ctx.HTML(http.StatusOK, template.DataSubsection, content)
	}
}

func (ds *DataSubsection[E]) CreateEntity(ctx *gin.Context) {
	var e E
	if err := ctx.ShouldBindJSON(&e); err != nil {
		ctx.HTML(http.StatusBadRequest, template.BadRequest, "")
		log.Println(err)
		return
	}

	if err := ds.uc.Create(&e); err != nil {
		ctx.HTML(http.StatusUnprocessableEntity, template.ErrDialog, gin.H{"message": err.Error()})
	} else {
		ds.sendUpdatedTable(ctx)
	}
}

func (ds *DataSubsection[E]) UpdateEntity(ctx *gin.Context) {
	var e E
	if err := ctx.ShouldBindJSON(&e); err != nil {
		ctx.HTML(http.StatusBadRequest, template.BadRequest, "")
		log.Println(err)
		return
	}

	if err := ds.uc.Update(&e); err != nil {
		ctx.HTML(http.StatusUnprocessableEntity, template.ErrDialog, gin.H{"message": err.Error()})
	} else {
		ds.sendUpdatedTable(ctx)
	}
}

func (ds *DataSubsection[E]) DeleteEntity(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.HTML(http.StatusBadRequest, template.BadRequest, "")
		log.Println("")
		return
	}

	if err := ds.uc.Delete(id); err != nil {
		ctx.HTML(http.StatusUnprocessableEntity, template.ErrDialog, gin.H{"message": err.Error()})
	} else {
		ds.sendUpdatedTable(ctx)
	}
}

func (ds *DataSubsection[E]) ExportData(ctx *gin.Context) {
	data, err := ds.uc.Export()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, template.InternalServerError, "")
		return
	}

	ctx.Header("Content-Type", "text/csv")

	wr := csv.NewWriter(ctx.Writer)
	if err := wr.WriteAll(data); err != nil {
		ctx.HTML(http.StatusInternalServerError, template.InternalServerError, "")
		return
	}
}

func (ds *DataSubsection[E]) ImportData(ctx *gin.Context) {

}

func (d *DataSubsection[E]) sendUpdatedTable(ctx *gin.Context) {
	if entities, err := d.uc.GetAll(); err != nil {
		ctx.HTML(http.StatusInternalServerError, template.InternalServerError, "")
	} else {
		content := d.prepareTableContent(entities)
		ctx.HTML(http.StatusOK, template.Table, content)
	}
}

func (d *DataSubsection[E]) preparePageContent(entities []E) *gin.H {
	var e E
	var endpoint, title string
	switch any(e).(type) {
	case entity.Peer:
		endpoint = PeersEndpoint
		title = "Peers"
	case entity.Friends:
		endpoint = FriendsEndpoint
		title = "Friends"
	case entity.Recommendations:
		endpoint = RecommendationsEndpoint
		title = "Recommendations"
	case entity.Task:
		endpoint = TasksEndpoint
		title = "Tasks"
	case entity.XP:
		endpoint = XPEndpont
		title = "XP"
	case entity.Points:
		endpoint = PointsEndpoint
		title = "Transferred points"
	case entity.Check:
		endpoint = ChecksEndpoint
		title = "Checks"
	case entity.P2P:
		endpoint = P2PEndpoint
		title = "P2P checks"
	case entity.Verter:
		endpoint = VerterEndpoint
		title = "Verter's checks"
	case entity.TimeTracking:
		endpoint = TimeTrackingEndpoint
		title = "Time tracking"
	}

	return &gin.H{
		"endpoint":      endpoint,
		"title":         title,
		"table_data":    entities,
		"table_headers": tools.GetFieldNames(e),
		"table_type":    "data",
	}
}

func (d *DataSubsection[E]) prepareTableContent(entities []E) (ginMap *gin.H) {
	var e E
	return &gin.H{
		"table_data":    entities,
		"table_headers": tools.GetFieldNames(e),
		"table_type":    "data",
	}
}
