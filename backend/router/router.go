package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zakirab-git/stockflow/config"
	"github.com/zakirab-git/stockflow/handler"
)

func SetupRouter(cfg config.Config, h *handler.Handlers) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	root := r.Group("/api")
	{
		products := root.Group("/products")
		{
			products.GET("", h.Product.GetProducts)
			products.POST("", h.Product.AddProduct)
			products.PUT("/:id", h.Product.UpdateProduct)
			products.DELETE("/:id", h.Product.DeleteProduct)
			products.POST("/:id/stock", h.Product.UpdateStock)
			products.GET("/:id/movements", h.Product.GetProductMovement)
		}

		orders := root.Group("/orders")
		{
			orders.POST("", h.Order.CreateOrder)
		}
	}

	return r
}
