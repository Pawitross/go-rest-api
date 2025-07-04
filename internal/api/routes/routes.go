package routes

import (
	"github.com/gin-gonic/gin"
	filesswag "github.com/swaggo/files"
	ginswag "github.com/swaggo/gin-swagger"

	_ "pawrest/docs"
	"pawrest/internal/api/handler"
	"pawrest/internal/api/middleware"
	"pawrest/internal/db"
)

// @title						Book managing API
// @description				Documentation of a book managing REST API.
// @BasePath					/api/v1
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and JWT token.

// @externalDocs.description	OpenAPI Specification
// @externalDocs.url			https://swagger.io/resources/open-api/
func Router(router *gin.Engine, db db.DatabaseInterface) {
	h := handler.Handlers{DB: db}

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			books := v1.Group("/books", middleware.Authenticate())
			{
				books.GET("", h.GetBooks)
				books.GET("/:id", h.GetBook)

				admin := books.Group("", middleware.Authorize())
				{
					admin.POST("", h.PostBook)

					admin.PUT("/:id", h.PutBook)

					admin.PATCH("/:id", h.PatchBook)

					admin.DELETE("/:id", h.DeleteBook)
				}
			}

			authors := v1.Group("/authors")
			{
				authors.GET("", h.GetAuthors)
				authors.GET("/:id", h.GetAuthor)

				authors.POST("", h.PostAuthor)

				authors.PUT("/:id", h.PutAuthor)

				authors.PATCH("/:id", h.PatchAuthor)

				authors.DELETE("/:id", h.DeleteAuthor)
			}

			v1.POST("login", handler.ReturnToken)
		}
	}

	router.GET("/swagger/*any", ginswag.WrapHandler(filesswag.Handler))
}
