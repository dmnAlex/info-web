package http

import (
	"log"
	"net/http"

	"github.com/dsnikitin/info-web/internal/entity"
	"github.com/dsnikitin/info-web/internal/pkg/tools"
	"github.com/gin-gonic/gin"
)

type DataSubsectionUseCase[E entity.Entity] interface {
	GetAll() ([]E, error)
	Create(e *E) error
	Update(e *E) error
	Delete(id string) error
}

type DataSubsection[E entity.Entity] struct {
	uc DataSubsectionUseCase[E]
}

func NewDataSubsection[E entity.Entity](uc DataSubsectionUseCase[E]) *DataSubsection[E] {
	return &DataSubsection[E]{uc: uc}
}

func (d *DataSubsection[E]) GetAllEntities(ctx *gin.Context) {
	if entities, err := d.uc.GetAll(); err != nil {
		ctx.HTML(http.StatusAlreadyReported, errDialogTemplate, gin.H{"message": err.Error()})
	} else {
		content := d.preparePageContent(entities)
		ctx.HTML(http.StatusOK, subsectionTemplate, content)
	}
}

func (d *DataSubsection[E]) CreateEntity(ctx *gin.Context) {
	var e E
	if err := ctx.Bind(&e); err != nil {
		// todo вернуть страницу bad request
		log.Println(err)
		return
	}

	if err := d.uc.Create(&e); err != nil {
		ctx.HTML(http.StatusAlreadyReported, errDialogTemplate, gin.H{"message": err.Error()})
	} else {
		d.sendUpdatedTable(ctx)
	}
}

func (d *DataSubsection[E]) UpdateEntity(ctx *gin.Context) {
	var e E
	if err := ctx.Bind(&e); err != nil {
		// todo вернуть страницу bad request
		log.Println(err)
		return
	}

	if err := d.uc.Update(&e); err != nil {
		ctx.HTML(http.StatusAlreadyReported, errDialogTemplate, gin.H{"message": err.Error()})
	} else {
		d.sendUpdatedTable(ctx)
	}
}

func (d *DataSubsection[E]) DeleteEntity(ctx *gin.Context) {
	id := ctx.Query("id")
	if err := d.uc.Delete(id); err != nil {
		ctx.HTML(http.StatusAlreadyReported, errDialogTemplate, gin.H{"message": err.Error()})
	} else {
		d.sendUpdatedTable(ctx)
	}
}

func (d *DataSubsection[E]) sendUpdatedTable(ctx *gin.Context) {
	if entities, err := d.uc.GetAll(); err != nil {
		ctx.HTML(http.StatusAlreadyReported, errDialogTemplate, gin.H{"message": err.Error()})
	} else {
		content := d.prepareTableContent(entities)
		ctx.HTML(http.StatusOK, tableTemplate, content)
	}
}

func (d *DataSubsection[E]) preparePageContent(entities []E) *gin.H {
	var e E
	var endpoint, title string
	switch any(e).(type) {
	case entity.Peer:
		endpoint = peersEndpoint
		title = "Peers"
	case entity.Friends:
		endpoint = friendsEndpoint
		title = "Friends"
	case entity.Recommendations:
		endpoint = recommendationsEndpoint
		title = "Recommendations"
	case entity.Task:
		endpoint = tasksEndpoint
		title = "Tasks"
	case entity.XP:
		endpoint = xpEndpont
		title = "XP"
	case entity.Points:
		endpoint = pointsEndpoint
		title = "Transferred points"
	case entity.Check:
		endpoint = checksEndpoint
		title = "Checks"
	case entity.P2P:
		endpoint = p2pEndpoint
		title = "P2P checks"
	case entity.Verter:
		endpoint = verterEndpoint
		title = "Verter's checks"
	case entity.TimeTracking:
		endpoint = timeTrackingEndpoint
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
