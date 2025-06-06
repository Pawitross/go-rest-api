package routes

import (
	"github.com/gin-gonic/gin"
	"pawrest/internal/api/handler"
)

func Router() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("books", handler.GetBooks)
			v1.GET("books/:id", handler.GetBook)

			v1.POST("books", handler.PostBook)

			v1.PUT("books/:id", handler.PutBook)

			v1.DELETE("books/:id", handler.DeleteBook)
		}
	}

	return router
}
