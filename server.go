package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vosgaust/covid19/entries"
	"github.com/vosgaust/covid19/http"
)

func run(cfg config) {
	r := gin.Default()

	entriesRepository := entries.NewMySQL(cfg.MySQL)
	entriesService := entries.NewService(entriesRepository)
	controller := http.NewEntriesController(entriesService)

	controller.AddRoutes(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
