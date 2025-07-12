package routes

import (
	"github.com/gin-gonic/gin"
	filesswag "github.com/swaggo/files"
	ginswag "github.com/swaggo/gin-swagger"

	_ "pawrest/docs"
	"pawrest/internal/api/handler"
	"pawrest/internal/api/middleware"
	"pawrest/internal/db"
	"pawrest/internal/yamlconfig"
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
// @description	To use limit, use the `limit` query parameter, like this: `limit=10`
// @description	To use offset, you also need to provide a limit.
// @description	The order of the limit and offset parameters doesn't matter.
// @description	Examples: `offset=10&limit=50`, `limit=50&offset=10`

// @BasePath	/api/v1

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description					Provide the JWT token as a Bearer token in the format "Bearer <your_token_here>".
// @description					To get the token use the /login endpoint with necessary body.

// @externalDocs.description	OpenAPI Specification
// @externalDocs.url			https://swagger.io/resources/open-api/
func Router(router *gin.Engine, db db.DatabaseInterface, cfg *yamlconfig.Config) {
	h := handler.Handlers{DB: db}
	secret := cfg.Secret

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			books := v1.Group("/books", middleware.Authenticate(secret))
			{
				books.GET("", h.GetBooks)
				books.GET("/:id", h.GetBook)
				books.OPTIONS("", h.OptionsBooks)
				books.OPTIONS("/:id", h.OptionsBook)

				admin := books.Group("", middleware.Authorize())
				{
					admin.POST("", h.PostBook)
					admin.PUT("/:id", h.PutBook)
					admin.PATCH("/:id", h.PatchBook)
					admin.DELETE("/:id", h.DeleteBook)
				}
			}

			authors := v1.Group("/authors", middleware.Authenticate(secret))
			{
				authors.GET("", h.GetAuthors)
				authors.GET("/:id", h.GetAuthor)
				authors.OPTIONS("", h.OptionsAuthors)
				authors.OPTIONS("/:id", h.OptionsAuthor)

				admin := authors.Group("", middleware.Authorize())
				{
					admin.POST("", h.PostAuthor)
					admin.PUT("/:id", h.PutAuthor)
					admin.PATCH("/:id", h.PatchAuthor)
					admin.DELETE("/:id", h.DeleteAuthor)
				}
			}

			genres := v1.Group("/genres", middleware.Authenticate(secret))
			{
				genres.GET("", h.GetGenres)
				genres.GET("/:id", h.GetGenre)
				genres.OPTIONS("", h.OptionsGenres)
				genres.OPTIONS("/:id", h.OptionsGenre)

				admin := genres.Group("", middleware.Authorize())
				{
					admin.POST("", h.PostGenre)
					admin.PUT("/:id", h.PutGenre)
					admin.DELETE("/:id", h.DeleteGenre)
				}
			}

			languages := v1.Group("/languages")
			{
				languages.GET("", h.GetLanguages)
				languages.GET("/:id", h.GetLanguage)
				languages.OPTIONS("", h.OptionsLanguages)
				languages.OPTIONS("/:id", h.OptionsLanguage)
				languages.POST("", h.PostLanguage)
				languages.PUT("/:id", h.PutLanguage)
				languages.DELETE("/:id", h.DeleteLanguage)
			}

			v1.POST("login", handler.ReturnToken(secret))
		}
	}

	router.GET("/swagger/*any", ginswag.WrapHandler(filesswag.Handler))
}
