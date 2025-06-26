package handler_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"pawrest/internal/api/routes"
	"pawrest/internal/db"
	"pawrest/internal/db/mock"
	"pawrest/internal/models"
)

var database db.BookDatabaseInterface

func runMain(m *testing.M) (int, error) {
	flag.Parse()

	if testing.Short() {
		database = &mock.MockDatabase{}
	} else {
		var err error
		database, err = db.ConnectToDB()
		if err != nil {
			return 0, err
		}
		defer database.(*db.Database).CloseDB()
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

func setupTestRouter(db db.BookDatabaseInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	routes.Router(router, db)

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

// GET /books
func TestListBooksSuccess(t *testing.T) {
	w := execRequest("GET", "/api/v1/books", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var rBooks []models.Book
	decodeBodyCheckEmpty(w, t, &rBooks)
}

func TestListBooksExtSuccess(t *testing.T) {
	w := execRequest("GET", "/api/v1/books?extend=true", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var rBooks []models.BookExt
	decodeBodyCheckEmpty(w, t, &rBooks)
}

func TestListBooksBadRequestUnknownParam(t *testing.T) {
	w := execRequest("GET", "/api/v1/books?foo=bar", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	checkErrorBodyNotEmpty(w, t)
}

// GET /books/id
func TestGetBookSuccess(t *testing.T) {
	w := execRequest("GET", "/api/v1/books/1", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var rBook models.Book
	decodeBodyCheckEmpty(w, t, &rBook)

	assert.NotEmpty(t, rBook, "Book in the response body should not be empty")
	assert.NotEmpty(t, rBook.Title, "Title should not be empty")
	assert.NotEmpty(t, rBook.Year, "Year should not be empty")
	assert.NotEmpty(t, rBook.Pages, "Pages should not be empty")
	assert.NotEmpty(t, rBook.Author, "Author should not be empty")
	assert.NotEmpty(t, rBook.Genre, "Genre should not be empty")
	assert.NotEmpty(t, rBook.Language, "Language should not be empty")
}

func TestGetBookNotFoundBigPathId(t *testing.T) {
	w := execRequest("GET", "/api/v1/books/9999", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)

	checkErrorBodyNotEmpty(w, t)
}

func TestGetBookBadRequestStringPathId(t *testing.T) {
	w := execRequest("GET", "/api/v1/books/string", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	checkErrorBodyNotEmpty(w, t)
}

// POST /books
func TestPostBookSuccess(t *testing.T) {
	testBook := models.Book{
		Title:    "Post book test",
		Year:     1996,
		Pages:    200,
		Author:   1,
		Genre:    1,
		Language: 1,
	}

	jsonBook := marshalCheckNoError(t, testBook)

	w := execRequest("POST", "/api/v1/books", bytes.NewReader(jsonBook))
	assert.Equal(t, http.StatusCreated, w.Code)

	var rBook models.Book
	decodeBodyCheckEmpty(w, t, &rBook)
	defer database.DelBook(rBook.Id)

	expLoc := fmt.Sprintf("/api/v1/books/%v", rBook.Id)
	assert.Equal(t, expLoc, w.Result().Header.Get("Location"))

	assert.NotZero(t, rBook.Id, "Auto generatated, non zero ID")
	assert.Equal(t, testBook.Title, rBook.Title)
	assert.Equal(t, testBook.Year, rBook.Year)
	assert.Equal(t, testBook.Pages, rBook.Pages)
	assert.Equal(t, testBook.Author, rBook.Author)
	assert.Equal(t, testBook.Genre, rBook.Genre)
	assert.Equal(t, testBook.Language, rBook.Language)
}

func TestPostBookBadRequestMalformedJSON(t *testing.T) {
	jsonStr := []byte(`{"title":"JSON Test","year":1996,"pages":200,"author":"Should be number","genre":1}`)
	w := execRequest("POST", "/api/v1/books", bytes.NewReader(jsonStr))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	checkErrorBodyNotEmpty(w, t)
}

func TestPostBookError(t *testing.T) {
	postTests := map[string]struct {
		testBook models.Book
		status   int
	}{
		"BadRequestValidationErr": {
			models.Book{
				Title:    "",
				Year:     1996,
				Pages:    200,
				Author:   1,
				Genre:    1,
				Language: 1,
			},
			http.StatusBadRequest,
		},
		"BadRequestForeignKeyErr": {
			models.Book{
				Title:    "Post foreign key test",
				Year:     1996,
				Pages:    200,
				Author:   1,
				Genre:    1,
				Language: 999,
			},
			http.StatusBadRequest,
		},
	}

	for name, tt := range postTests {
		t.Run(name, func(t *testing.T) {
			jsonBook := marshalCheckNoError(t, tt.testBook)

			w := execRequest("POST", "/api/v1/books", bytes.NewReader(jsonBook))
			assert.Equal(t, tt.status, w.Code)

			checkErrorBodyNotEmpty(w, t)
		})
	}
}

// PUT /books/id
func TestPutBookSuccess(t *testing.T) {
	testBook := models.Book{
		Title:    "Put book test",
		Year:     1996,
		Pages:    593,
		Author:   1,
		Genre:    1,
		Language: 1,
	}

	jsonBook := marshalCheckNoError(t, testBook)

	w := execRequest("PUT", "/api/v1/books/1", bytes.NewReader(jsonBook))
	assert.Equal(t, http.StatusNoContent, w.Code)

	book, _ := database.GetBook(1)
	assert.Equal(t, testBook.Title, book.Title)
	assert.Equal(t, testBook.Year, book.Year)
	assert.Equal(t, testBook.Pages, book.Pages)
	assert.Equal(t, testBook.Author, book.Author)
	assert.Equal(t, testBook.Genre, book.Genre)
	assert.Equal(t, testBook.Language, book.Language)
}

func TestPutBookBadRequestMalformedJSON(t *testing.T) {
	jsonStr := []byte(`{"title":"JSON Test","year":1996,"pages":200,"author":"Should be number","genre":1}`)
	w := execRequest("PUT", "/api/v1/books/1", bytes.NewReader(jsonStr))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	checkErrorBodyNotEmpty(w, t)
}

func TestPutBookError(t *testing.T) {
	putTests := map[string]struct {
		testBook models.Book
		query    string
		status   int
	}{
		"NotFoundBigPathId": {
			models.Book{
				Title:    "Put book test",
				Year:     1996,
				Pages:    593,
				Author:   1,
				Genre:    1,
				Language: 1,
			},
			"/9999",
			http.StatusNotFound,
		},
		"BadRequestForeignKeyErr": {
			models.Book{
				Title:    "Put foreign key test",
				Year:     1996,
				Pages:    593,
				Author:   1,
				Genre:    1,
				Language: 999,
			},
			"/1",
			http.StatusBadRequest,
		},
		"BadRequestStringId": {
			models.Book{
				Title:    "Put book test",
				Year:     1996,
				Pages:    593,
				Author:   1,
				Genre:    1,
				Language: 1,
			},
			"/string",
			http.StatusBadRequest,
		},
		"BadRequestBadJSON": {
			models.Book{
				Title: "Put book test",
				Year:  1996,
				Pages: 593,
			},
			"/1",
			http.StatusBadRequest,
		},
	}

	for name, tt := range putTests {
		t.Run(name, func(t *testing.T) {
			fullUrl := "/api/v1/books" + tt.query

			jsonBook := marshalCheckNoError(t, tt.testBook)

			w := execRequest("PUT", fullUrl, bytes.NewReader(jsonBook))
			assert.Equal(t, tt.status, w.Code)

			checkErrorBodyNotEmpty(w, t)
		})
	}
}

// PATCH /books/id
func TestPatchBookSuccess(t *testing.T) {
	jsonStr := []byte(`{"title":"Patch book test", "pages":999}`)
	w := execRequest("PATCH", "/api/v1/books/1", bytes.NewReader(jsonStr))
	assert.Equal(t, http.StatusNoContent, w.Code)

	book, _ := database.GetBook(1)
	assert.Equal(t, "Patch book test", book.Title)
	assert.Equal(t, int64(999), book.Pages)
}

func TestPatchBookError(t *testing.T) {
	patchTests := map[string]struct {
		jsonStr string
		query   string
		status  int
	}{
		"NotFoundBigPathId": {
			`{"title":"Patch book test", "pages":999}`,
			"/9999",
			http.StatusNotFound,
		},
		"BadRequestStringPathId": {
			`{"title":"Patch book test", "pages":999}`,
			"/string",
			http.StatusBadRequest,
		},
		"BadRequestBadJSON": {
			`{"title":"Patch book test", "pages":Should be number"}`,
			"/1",
			http.StatusBadRequest,
		},
	}

	for name, tt := range patchTests {
		t.Run(name, func(t *testing.T) {
			fullUrl := "/api/v1/books" + tt.query

			w := execRequest("PATCH", fullUrl, bytes.NewReader([]byte(tt.jsonStr)))
			assert.Equal(t, tt.status, w.Code)

			checkErrorBodyNotEmpty(w, t)
		})
	}
}

// DELETE /books/id
func TestDeleteBookSuccess(t *testing.T) {
	if !testing.Short() {
		t.Skip("Skipping success deletion test on real database")
	}

	w := execRequest("DELETE", "/api/v1/books/2", nil)
	assert.Equal(t, http.StatusNoContent, w.Code)

	_, err := database.GetBook(2)
	assert.ErrorIs(t, err, db.ErrNotFound)
}

func TestDeleteBookError(t *testing.T) {
	deleteTests := map[string]struct {
		query  string
		status int
	}{
		"NotFoundBigPathId": {
			"/9999",
			http.StatusNotFound,
		},
		"BadRequestStringPathId": {
			"/string",
			http.StatusBadRequest,
		},
	}

	for name, tt := range deleteTests {
		t.Run(name, func(t *testing.T) {
			fullUrl := "/api/v1/books" + tt.query

			w := execRequest("DELETE", fullUrl, nil)
			assert.Equal(t, tt.status, w.Code)

			checkErrorBodyNotEmpty(w, t)
		})
	}
}
