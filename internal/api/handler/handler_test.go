package handler_test

import (
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
)

var database db.DatabaseInterface

func runMain(m *testing.M) (int, error) {
	flag.Parse()

	if testing.Short() {
		database = &mock.MockDatabase{}
	} else {
		os.Setenv("DBUSER", "root")
		os.Setenv("DBPASS", "")
		os.Setenv("DBNAME", "paw_test")

		var err error
		database, err = db.ConnectToDB()
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

			books.POST("", h.PostBook)

			books.PUT("/:id", h.PutBook)

			books.PATCH("/:id", h.PatchBook)

			books.DELETE("/:id", h.DeleteBook)
		}

		authors := apiv1.Group("/authors")
		{
			authors.GET("", h.GetAuthors)
			authors.GET("/:id", h.GetAuthor)

			authors.POST("", h.PostAuthor)

			authors.PUT("/:id", h.PutAuthor)

			authors.PATCH("/:id", h.PatchAuthor)

			authors.DELETE("/:id", h.DeleteAuthor)
		}

		apiv1.POST("login", handler.ReturnToken)
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

func marshalCheckNoError(t *testing.T, obj any) []byte {
	j, err := json.Marshal(obj)
	assert.NoError(t, err, "JSON marshalling error")

	return j
}

func decodeBodyCheckEmpty(w *httptest.ResponseRecorder, t *testing.T, obj any) {
	err := json.NewDecoder(w.Body).Decode(obj)
	assert.NoError(t, err, "Error decoding response data:", err)
	assert.NotEmpty(t, obj, "Obj should not be empty")
}

func checkErrorBodyNotEmpty(w *httptest.ResponseRecorder, t *testing.T) {
	var rError models.Error
	decodeBodyCheckEmpty(w, t, &rError)
}
