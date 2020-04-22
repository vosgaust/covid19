package http

import (
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/vosgaust/covid19/entries"
)

// Controller is used by gin routes to handle requests
type Controller struct {
	entriesService entries.Service
}

// NewEntriesController creates an instance of entries controller
func NewEntriesController(entriesService entries.Service) *Controller {
	return &Controller{entriesService}
}

// AddRoutes setup common configuration such as cors
func (c *Controller) AddRoutes(g *gin.Engine) {
	g.Use(setCors())

	v1 := g.Group("/api/v1")

	v1.GET("/historic/:state/cumulative", c.getCumulative())
	v1.GET("/historic/:state/deltas", c.getDeltas())
}

func setCors() gin.HandlerFunc {
	corsCfg := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	return corsCfg
}
