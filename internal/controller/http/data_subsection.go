package http

import (
	"log"
	"net/http"

	"github.com/dsnikitin/info-web/internal/entity"
	"github.com/gin-gonic/gin"
)

type DataSubsectionUseCase[E entity.Entity] interface {
	GetAll() []E
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
	entities := d.uc.GetAll()
	template, content := d.preparePageContent(entities)
	ctx.HTML(http.StatusOK, template, content)
}

func (d *DataSubsection[E]) CreateEntity(ctx *gin.Context) {
	var e E
	if err := ctx.Bind(&e); err != nil {
		// todo вернуть страницу bad request
		log.Println(err)
		return
	}

	if err := d.uc.Create(&e); err != nil {
		ctx.HTML(http.StatusAlreadyReported, errDialogTemplate, gin.H{"message": entity.GetErrorDescription(err)})
	} else {
		ctx.HTML(http.StatusOK, tableTemplate, d.prepareTableContent(d.uc.GetAll()))
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
		ctx.HTML(http.StatusAlreadyReported, errDialogTemplate, gin.H{"message": entity.GetErrorDescription(err)})
	} else {
		ctx.HTML(http.StatusOK, tableTemplate, d.prepareTableContent(d.uc.GetAll()))
	}
}

func (d *DataSubsection[E]) DeleteEntity(ctx *gin.Context) {
	id := ctx.Query("id")

	if err := d.uc.Delete(id); err != nil {
		ctx.HTML(http.StatusAlreadyReported, errDialogTemplate, gin.H{"message": entity.GetErrorDescription(err)})
	} else {
		ctx.HTML(http.StatusOK, tableTemplate, d.prepareTableContent(d.uc.GetAll()))
	}
}

func (d *DataSubsection[E]) preparePageContent(entities []E) (template string, ginMap *gin.H) {
	switch any(entities[0]).(type) {
	case entity.Peer:
		template = peersTemplate
		ginMap = &gin.H{
			"endpoint":      peersEndpoint,
			"object_fields": []string{"nickname", "birthday"},
			"table_title":   "Студенты",
			"table_data":    entities,
			"table_headers": []string{"ID", "Ник", "Дата рождения"},
			"table_type":    "data",
		}
	case entity.Friends:
		template = friendsTemplate
		ginMap = &gin.H{
			"endpoint":      friendsEndpoint,
			"object_fields": []string{"peer_1", "peer_2"},
			"table_title":   "Друзья",
			"table_data":    entities,
			"table_headers": []string{"ID", "Ник студента", "Ник друга"},
			"table_type":    "data",
		}
	case entity.Recommendations:
		template = recommendationsTemplate
		ginMap = &gin.H{
			"endpoint":      recommendationsEndpoint,
			"object_fields": []string{"peer", "recommended_peer"},
			"table_title":   "Рекоммендации",
			"table_data":    entities,
			"table_headers": []string{"ID", "Ник студента", "Ник рекоммендуемого"},
			"table_type":    "data",
		}
	case entity.Task:
		template = tasksTemplate
		ginMap = &gin.H{
			"endpoint":      tasksEndpoint,
			"object_fields": []string{"title", "parent_task", "max_xp"},
			"table_title":   "Задания",
			"table_data":    entities,
			"table_headers": []string{"ID", "Название", "Родительское задание", "Максимальный опыт"},
			"table_type":    "data",
		}
	case entity.XP:
		template = xpTemplate
		ginMap = &gin.H{
			"endpoint":      xpEndpont,
			"object_fields": []string{"check_id", "xp_amount"},
			"table_title":   "Опыт",
			"table_data":    entities,
			"table_headers": []string{"ID", "ID проверки", "Количество опыта"},
			"table_type":    "data",
		}
	case entity.Points:
		template = pointsTemplate
		ginMap = &gin.H{
			"endpoint":      pointsEndpoint,
			"object_fields": []string{"checking_peer", "checked_peer", "points_amount"},
			"table_title":   "Полученные поинты",
			"table_data":    entities,
			"table_headers": []string{"ID", "Проверяющий", "Проверяемый", "Количество поинтов"},
			"table_type":    "data",
		}
	case entity.Check:
		template = checksTemplate
		ginMap = &gin.H{
			"endpoint":      checksEndpoint,
			"object_fields": []string{"peer", "task", "date"},
			"table_title":   "Провероки",
			"table_data":    entities,
			"table_headers": []string{"ID", "Ник студента", "Название задания", "Дата проверки"},
			"table_type":    "data",
		}
	case entity.P2P:
		template = p2pTemplate
		ginMap = &gin.H{
			"endpoint":      p2pEndpoint,
			"object_fields": []string{"check_id", "checking_peer", "state", "time"},
			"table_title":   "Проверки p2p",
			"table_data":    entities,
			"table_headers": []string{"ID", "ID проверки", "Ник проверяющего", "Статус проверки", "Время"},
			"table_type":    "data",
		}
	case entity.Verter:
		template = verterTemplate
		ginMap = &gin.H{
			"endpoint":      verterEndpoint,
			"object_fields": []string{"checking_id", "state", "time"},
			"table_title":   "Проверок вертера",
			"table_data":    entities,
			"table_headers": []string{"ID", "ID проверки", "Статус проверки", "Время"},
			"table_type":    "data",
		}
	case entity.TimeTracking:
		template = timeTrackingTemplate
		ginMap = &gin.H{
			"endpoint":      timeTrackingEndpoint,
			"object_fields": []string{"peer", "date", "time", "state"},
			"table_title":   "Посещения",
			"table_data":    entities,
			"table_headers": []string{"ID", "Ник студента", "Дата", "Время", "Статус"},
			"table_type":    "data",
		}
	}

	return template, ginMap
}

func (d *DataSubsection[E]) prepareTableContent(entities []E) (ginMap *gin.H) {
	switch any(entities[0]).(type) {
	case entity.Peer:
		ginMap = &gin.H{
			"table_data":    entities,
			"table_headers": []string{"ID", "Ник", "Дата рождения"},
			"table_type":    "data",
		}
	case entity.Friends:
		ginMap = &gin.H{
			"table_data":    entities,
			"table_headers": []string{"ID", "Ник студента", "Ник друга"},
			"table_type":    "data",
		}
	case entity.Recommendations:
		ginMap = &gin.H{
			"table_data":    entities,
			"table_headers": []string{"ID", "Ник студента", "Ник рекоммендуемого"},
			"table_type":    "data",
		}
	case entity.Task:
		ginMap = &gin.H{
			"table_data":    entities,
			"table_headers": []string{"ID", "Название", "Родительское задание", "Максимальный опыт"},
			"table_type":    "data",
		}
	case entity.XP:
		ginMap = &gin.H{
			"table_data":    entities,
			"table_headers": []string{"ID", "ID проверки", "Количество опыта"},
			"table_type":    "data",
		}
	case entity.Points:
		ginMap = &gin.H{
			"table_data":    entities,
			"table_headers": []string{"ID", "Проверяющий", "Проверяемый", "Количество поинтов"},
			"table_type":    "data",
		}
	case entity.Check:
		ginMap = &gin.H{
			"table_data":    entities,
			"table_headers": []string{"ID", "Ник студента", "Название задания", "Дата проверки"},
			"table_type":    "data",
		}
	case entity.P2P:
		ginMap = &gin.H{
			"table_data":    entities,
			"table_headers": []string{"ID", "ID проверки", "Ник проверяющего", "Статус проверки", "Время"},
			"table_type":    "data",
		}
	case entity.Verter:
		ginMap = &gin.H{
			"table_data":    entities,
			"table_headers": []string{"ID", "ID проверки", "Статус проверки", "Время"},
			"table_type":    "data",
		}
	case entity.TimeTracking:
		ginMap = &gin.H{
			"table_data":    entities,
			"table_headers": []string{"ID", "Ник студента", "Дата", "Время", "Статус"},
			"table_type":    "data",
		}
	}

	return ginMap
}
