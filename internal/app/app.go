package app

import (
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/dsnikitin/info-web/internal/pkg/tools"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Run() {
	db, err := gorm.Open(postgres.Open("sslmode=disable"), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatal("Failed to connect to db")
	}
	log.Println("Established a successful connection to db!")

	repos := initRepositories(db)
	usecases := initUseCases(repos)
	handlers := initHandlers(usecases)

	router := gin.Default()
	router.SetFuncMap(template.FuncMap{ // must be invoked before LoadHTMLFiles
		"add":                tools.Add,
		"getFieldValues":     tools.GetFieldValues,
		"getPrimaryKeyValue": tools.GetPrimaryKeyValue,
		"toLowerCase":        tools.ToLowerCase,
	})
	router.LoadHTMLFiles(tools.GetAllTemplates()...)
	router.Static("/internal/assets", "internal/assets")
	initRoutes(router, handlers)

	router.Run(fmt.Sprintf(":%v", os.Getenv("APPPORT")))
}
