package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakirab-git/stockflow/config"
	"github.com/zakirab-git/stockflow/handler"
)

func SetupRouter(cfg config.Config, h *handler.Handlers) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	root := r.Group("/api/v1")
	{
		test := root.Group("/test")
		{
			test.GET("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "pong"})
			})
		}
	}

	return r
}
