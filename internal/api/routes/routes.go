package routes

import (
	"github.com/gin-gonic/gin"
	filesswag "github.com/swaggo/files"
	ginswag "github.com/swaggo/gin-swagger"

	"pawrest/docs"
	"pawrest/internal/api/handler"
	mware "pawrest/internal/api/middleware"
)

func Router(router *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/api/v1"
	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			books := v1.Group("/books")
			{
				books.GET("", handler.GetBooks)
				books.GET("/:id", handler.GetBook)

				books.POST("", handler.PostBook)

				books.PUT("/:id", handler.PutBook)

				books.PATCH("/:id", handler.PatchBook)

				books.DELETE("/:id", handler.DeleteBook)
			}

			v1.GET("login", handler.ReturnToken)

			v1.GET("auth", mware.Authenticate(), mware.Authorize(), func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Witamy"})
			})
		}
	}

	router.GET("/swagger/*any", ginswag.WrapHandler(filesswag.Handler))
}
