package app

import (
	"log"
	"os"
	"path/filepath"
	"reflect"

	"html/template"

	"github.com/dsnikitin/info-web/internal/controller/http"
	"github.com/dsnikitin/info-web/internal/entity"
	"github.com/dsnikitin/info-web/internal/infrastructure/repository"
	"github.com/dsnikitin/info-web/internal/usecase"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repositories struct {
	operations      *repository.Operations
	checks          *repository.DataManager[entity.Check]
	friends         *repository.DataManager[entity.Friends]
	p2p             *repository.DataManager[entity.P2P]
	peers           *repository.DataManager[entity.Peer]
	recommendations *repository.DataManager[entity.Recommendations]
	tasks           *repository.DataManager[entity.Task]
	time_tracking   *repository.DataManager[entity.TimeTracking]
	points          *repository.DataManager[entity.Points]
	verter          *repository.DataManager[entity.Verter]
	xp              *repository.DataManager[entity.XP]
}

type UseCases struct {
	operations      *usecase.Operations
	checks          *usecase.DataSubsection[entity.Check]
	friends         *usecase.DataSubsection[entity.Friends]
	p2p             *usecase.DataSubsection[entity.P2P]
	peers           *usecase.DataSubsection[entity.Peer]
	recommendations *usecase.DataSubsection[entity.Recommendations]
	tasks           *usecase.DataSubsection[entity.Task]
	time_tracking   *usecase.DataSubsection[entity.TimeTracking]
	points          *usecase.DataSubsection[entity.Points]
	verter          *usecase.DataSubsection[entity.Verter]
	xp              *usecase.DataSubsection[entity.XP]
}

type Handlers struct {
	index           *http.Index
	data            *http.Data
	operations      *http.Operations
	checks          *http.DataSubsection[entity.Check]
	friends         *http.DataSubsection[entity.Friends]
	p2p             *http.DataSubsection[entity.P2P]
	peers           *http.DataSubsection[entity.Peer]
	recommendations *http.DataSubsection[entity.Recommendations]
	tasks           *http.DataSubsection[entity.Task]
	time_tracking   *http.DataSubsection[entity.TimeTracking]
	points          *http.DataSubsection[entity.Points]
	verter          *http.DataSubsection[entity.Verter]
	xp              *http.DataSubsection[entity.XP]
}

func Run() {
	db, err := ConnectToDB()
	if err != nil {
		log.Fatal("Failed to connect to db")
	}
	log.Println("Established a successful connection to db!")

	repos := initRepositories(db)
	usecases := initUseCases(repos)
	handlers := initHandlers(usecases)

	engine := gin.Default()
	initTemplates(engine)
	initRoutes(engine, handlers)

	engine.Run(":9060")
}

func initRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		operations:      repository.NewOperations(db),
		checks:          repository.NewDataManager[entity.Check](db),
		friends:         repository.NewDataManager[entity.Friends](db),
		p2p:             repository.NewDataManager[entity.P2P](db),
		peers:           repository.NewDataManager[entity.Peer](db),
		recommendations: repository.NewDataManager[entity.Recommendations](db),
		tasks:           repository.NewDataManager[entity.Task](db),
		time_tracking:   repository.NewDataManager[entity.TimeTracking](db),
		points:          repository.NewDataManager[entity.Points](db),
		verter:          repository.NewDataManager[entity.Verter](db),
		xp:              repository.NewDataManager[entity.XP](db),
	}
}

func ConnectToDB() (*gorm.DB, error) {
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=postgres"
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func initUseCases(rs *Repositories) *UseCases {
	return &UseCases{
		operations:      usecase.NewOperations(rs.operations),
		checks:          usecase.NewDataSubsection[entity.Check](rs.checks),
		friends:         usecase.NewDataSubsection[entity.Friends](rs.friends),
		p2p:             usecase.NewDataSubsection[entity.P2P](rs.p2p),
		peers:           usecase.NewDataSubsection[entity.Peer](rs.peers),
		recommendations: usecase.NewDataSubsection[entity.Recommendations](rs.recommendations),
		tasks:           usecase.NewDataSubsection[entity.Task](rs.tasks),
		time_tracking:   usecase.NewDataSubsection[entity.TimeTracking](rs.time_tracking),
		points:          usecase.NewDataSubsection[entity.Points](rs.points),
		verter:          usecase.NewDataSubsection[entity.Verter](rs.verter),
		xp:              usecase.NewDataSubsection[entity.XP](rs.xp),
	}
}

func initHandlers(uc *UseCases) *Handlers {
	return &Handlers{
		index:           http.NewIndex(),
		data:            http.NewData(),
		operations:      http.NewOperations(uc.operations),
		checks:          http.NewDataSubsection[entity.Check](uc.checks),
		friends:         http.NewDataSubsection[entity.Friends](uc.friends),
		p2p:             http.NewDataSubsection[entity.P2P](uc.p2p),
		peers:           http.NewDataSubsection[entity.Peer](uc.peers),
		recommendations: http.NewDataSubsection[entity.Recommendations](uc.recommendations),
		tasks:           http.NewDataSubsection[entity.Task](uc.tasks),
		time_tracking:   http.NewDataSubsection[entity.TimeTracking](uc.time_tracking),
		points:          http.NewDataSubsection[entity.Points](uc.points),
		verter:          http.NewDataSubsection[entity.Verter](uc.verter),
		xp:              http.NewDataSubsection[entity.XP](uc.xp),
	}
}

func initTemplates(gin *gin.Engine) {
	templates := make([]string, 0)
	err := filepath.Walk("internal/templates/", func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			templates = append(templates, path)
		}
		return err
	})
	if err != nil {
		log.Fatalf("%v", err)
	}

	gin.SetFuncMap(template.FuncMap{
		"add":            add,
		"getFieldsValue": getFieldsValue,
	})
	gin.LoadHTMLFiles(templates...)
	gin.Static("/internal/assets", "internal/assets")

}

func add(a int, b int) int {
	return a + b
}

// returns the value of each field of struct in a slice
func getFieldsValue(s interface{}) []interface{} {
	if s == nil {
		return nil
	}

	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Struct {
		return nil
	}

	out := make([]interface{}, 0, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		out = append(out, v.Field(i).Interface())
	}

	return out
}

func initRoutes(e *gin.Engine, h *Handlers) {
	e.GET("/", h.index.GetIndexPage)
	e.GET("/data", h.data.GetFrontPage)
	e.GET("/operations", h.operations.GetAll)

	e.GET("/data/checks", h.checks.GetAllEntities)
	e.POST("/data/checks", h.checks.CreateEntity)
	e.PUT("/data/checks", h.checks.UpdateEntity)
	e.DELETE("/data/checks", h.checks.DeleteEntity)

	e.GET("/data/friends", h.friends.GetAllEntities)
	e.POST("/data/friends", h.friends.CreateEntity)
	e.PUT("/data/friends", h.friends.UpdateEntity)
	e.DELETE("/data/friends", h.friends.DeleteEntity)

	e.GET("/data/p2p", h.p2p.GetAllEntities)
	e.POST("/data/p2p", h.p2p.CreateEntity)
	e.PUT("/data/p2p", h.p2p.UpdateEntity)
	e.DELETE("/data/p2p", h.p2p.DeleteEntity)

	e.GET("/data/peers", h.peers.GetAllEntities)
	e.POST("/data/peers", h.peers.CreateEntity)
	e.PUT("/data/peers", h.peers.UpdateEntity)
	e.DELETE("/data/peers", h.peers.DeleteEntity)

	e.GET("/data/recommendations", h.recommendations.GetAllEntities)
	e.POST("/data/recommendations", h.recommendations.CreateEntity)
	e.PUT("/data/recommendations", h.recommendations.UpdateEntity)
	e.DELETE("/data/recommendations", h.recommendations.DeleteEntity)

	e.GET("/data/tasks", h.tasks.GetAllEntities)
	e.POST("/data/tasks", h.tasks.CreateEntity)
	e.PUT("/data/tasks", h.tasks.UpdateEntity)
	e.DELETE("/data/tasks", h.tasks.DeleteEntity)

	e.GET("/data/time_tracking", h.time_tracking.GetAllEntities)
	e.POST("/data/time_tracking", h.time_tracking.CreateEntity)
	e.PUT("/data/time_tracking", h.time_tracking.UpdateEntity)
	e.DELETE("/data/time_tracking", h.time_tracking.DeleteEntity)

	e.GET("/data/points", h.points.GetAllEntities)
	e.POST("/data/points", h.points.CreateEntity)
	e.PUT("/data/points", h.points.UpdateEntity)
	e.DELETE("/data/tpoints", h.points.DeleteEntity)

	e.GET("/data/verter", h.verter.GetAllEntities)
	e.POST("/data/verter", h.verter.CreateEntity)
	e.PUT("/data/verter", h.verter.UpdateEntity)
	e.DELETE("/data/verter", h.verter.DeleteEntity)

	e.GET("/data/xp", h.xp.GetAllEntities)
	e.POST("/data/xp", h.xp.CreateEntity)
	e.PUT("/data/xp", h.xp.UpdateEntity)
	e.DELETE("/data/xp", h.xp.DeleteEntity)
}
