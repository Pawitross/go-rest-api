package routes

import (
	"github.com/gin-gonic/gin"
	"pawrest/internal/api/handler"
	mware "pawrest/internal/api/middleware"
)

func Router(router *gin.Engine) {
	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("books", handler.GetBooks)
			v1.GET("books/:id", handler.GetBook)

			v1.POST("books", handler.PostBook)

			v1.PUT("books/:id", handler.PutBook)

			v1.PATCH("books/:id", handler.PatchBook)

			v1.DELETE("books/:id", handler.DeleteBook)

			v1.GET("login", handler.ReturnToken)

			v1.GET("auth", mware.Authenticate(), mware.Authorize(), func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Witamy"})
			})
		}
	}
}
