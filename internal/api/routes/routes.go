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

// @title		Book managing API
// @description	Documentation of a book managing REST API.
// @description
// @description	**How to use filtering:**
// @description	To use simple filtering put name of the column in the query parameter followed by the value.
// @description	Examples: `last_name=Orwell`, `title=Dziady`
// @description	To filter extended response use filtering like this: `genre.name=Nowela`
// @description
// @description	To filter using comparison operators append the operator to the query parameter. Available operators:
// @description	- less than = `.lt`
// @description	- less than or equal = `.lte`
// @description	- greater than = `.gt`
// @description	- greater than or equal = `.gte`
// @description	- equal = `.eq`
// @description	- not equal = `.neq`
// @description
// @description	Examples: `pages.lt=300`, `year.gte=1980`, `language.name.neq=Polski`.
// @description
// @description	**How to use sorting:**
// @description	To sort, use `sort_by` query parameter followed by the column name.
// @description	If you want to sort in descending order, prefix the column name with a minus sign (`-`).
// @description	Examples: `sort_by=pages` - ascending order, `sort_by=-pages` - descending order
// @description
// @description	**How to use limit and offset:**
// @description	To use a limit, use the `limit` query parameter, like this: `limit=10`
// @description	To use offset, you also need to provide a limit as well.
// @description	The order of the limit and offset doesn't matter.
// @description	Examples: `offset=10&limit=50`, `limit=50&offset=10`

// @BasePath					/api/v1

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description					Provide the JWT token as a Bearer token in the format "Bearer <your_token_here>".
// @description					To get the token use the /login endpoint with necessary body.

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
