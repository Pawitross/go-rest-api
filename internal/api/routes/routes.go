package routes

import (
	"github.com/gin-gonic/gin"
	filesswag "github.com/swaggo/files"
	ginswag "github.com/swaggo/gin-swagger"

	"pawrest/docs"
	"pawrest/internal/api/handler"
	"pawrest/internal/api/middleware"
	"pawrest/internal/db"
)

func Router(router *gin.Engine, db *db.Database) {
	h := handler.Handlers{DB: db}

	docs.SwaggerInfo.BasePath = "/api/v1"
	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			books := v1.Group("/books")
			{
				books.GET("", h.GetBooks)
				books.GET("/:id", h.GetBook)

				books.POST("", h.PostBook)

				books.PUT("/:id", h.PutBook)

				books.PATCH("/:id", h.PatchBook)

				books.DELETE("/:id", h.DeleteBook)
			}

			v1.GET("login", handler.ReturnToken)

			v1.GET("auth", middleware.Authenticate(), middleware.Authorize(), func(c *gin.Context) {
				c.JSON(200, "Welcome")
			})
		}
	}

	router.GET("/swagger/*any", ginswag.WrapHandler(filesswag.Handler))
}
