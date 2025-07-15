package handler_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"pawrest/internal/api/handler"
	"pawrest/internal/db"
	"pawrest/internal/db/mock"
	"pawrest/internal/models"
	"pawrest/internal/testutil"
	"pawrest/internal/yamlconfig"
)

var database db.DatabaseInterface

func runMain(m *testing.M) (int, error) {
	flag.Parse()

	if testing.Short() {
		database = mock.NewMockDatabase()
	} else {
		cfg := &yamlconfig.Config{
			DBUser: "user_test",
			DBPass: "testpass",
			DBName: "paw_test",
			DBHost: "127.0.0.1",
			DBPort: "3306",
		}

		var err error
		database, err = db.ConnectToDB(cfg)
		if err != nil {
			return 0, err
		}
		defer database.(*db.Database).CloseDB()

		if err := testutil.SetupDatabase(database.(*db.Database).Pool()); err != nil {
			return 0, err
		}
	}

	return m.Run(), nil
}

func TestMain(m *testing.M) {
	code, err := runMain(m)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

func execAndCheck(t *testing.T, method, url string, body []byte, status int, o any) *httptest.ResponseRecorder {
	t.Helper()
	w := execRequest(method, url, bytes.NewReader(body))
	assert.Equal(t, status, w.Code)

	if o != nil {
		decodeJSONBodyCheckEmpty(t, w, o)
	}

	return w
}

func execAndCheckError(t *testing.T, method, url string, body []byte, status int) {
	var rError models.Error
	execAndCheck(t, method, url, body, status, &rError)
}

func setupTestRouter(db db.DatabaseInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	h := handler.Handlers{DB: db}

	apiv1 := router.Group("/api/v1")
	{
		books := apiv1.Group("/books")
		{
			books.GET("", h.GetBooks)
			books.GET("/:id", h.GetBook)
			books.OPTIONS("", h.OptionsBooks)
			books.OPTIONS("/:id", h.OptionsBook)
			books.POST("", h.PostBook)
			books.PUT("/:id", h.PutBook)
			books.PATCH("/:id", h.PatchBook)
			books.DELETE("/:id", h.DeleteBook)
		}

		authors := apiv1.Group("/authors")
		{
			authors.GET("", h.GetAuthors)
			authors.GET("/:id", h.GetAuthor)
			authors.OPTIONS("", h.OptionsAuthors)
			authors.OPTIONS("/:id", h.OptionsAuthor)
			authors.POST("", h.PostAuthor)
			authors.PUT("/:id", h.PutAuthor)
			authors.PATCH("/:id", h.PatchAuthor)
			authors.DELETE("/:id", h.DeleteAuthor)
		}

		genres := apiv1.Group("/genres")
		{
			genres.GET("", h.GetGenres)
			genres.GET("/:id", h.GetGenre)
			genres.OPTIONS("", h.OptionsGenres)
			genres.OPTIONS("/:id", h.OptionsGenre)
			genres.POST("", h.PostGenre)
			genres.PUT("/:id", h.PutGenre)
			genres.DELETE("/:id", h.DeleteGenre)
		}

		languages := apiv1.Group("/languages")
		{
			languages.GET("", h.GetLanguages)
			languages.GET("/:id", h.GetLanguage)
			languages.OPTIONS("", h.OptionsLanguages)
			languages.OPTIONS("/:id", h.OptionsLanguage)
			languages.POST("", h.PostLanguage)
			languages.PUT("/:id", h.PutLanguage)
			languages.DELETE("/:id", h.DeleteLanguage)
		}

		apiv1.POST("login", handler.ReturnToken("random-string"))
	}

	return router
}

func execRequest(method, target string, body io.Reader) *httptest.ResponseRecorder {
	router := setupTestRouter(database)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, target, body)

	router.ServeHTTP(w, req)
	return w
}

func decodeJSONBodyCheckEmpty(t *testing.T, w *httptest.ResponseRecorder, obj any) {
	t.Helper()
	err := json.NewDecoder(w.Body).Decode(obj)
	assert.NoError(t, err, "Error decoding response data:", err)
	assert.NotEmpty(t, obj, "Obj should not be empty")
}

func marshalCheckNoError(t *testing.T, obj any) []byte {
	t.Helper()
	j, err := json.Marshal(obj)
	assert.NoError(t, err, "JSON marshalling error")

	return j
}
