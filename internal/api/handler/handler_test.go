package handler_test

import (
	"bytes"
	"encoding/json"
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
	"pawrest/internal/models"
)

var database *db.Database

func runMain(m *testing.M) (int, error) {
	var err error
	database, err = db.ConnectToDB()
	if err != nil {
		return 0, err
	}
	defer database.CloseDB()

	return m.Run(), nil
}

func TestMain(m *testing.M) {
	code, err := runMain(m)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

func setupTestRouter(db *db.Database) *gin.Engine {
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

func TestPostBookBadRequestValidationErr(t *testing.T) {
	testBook := models.Book{
		Title:    "",
		Year:     1996,
		Pages:    200,
		Author:   1,
		Genre:    1,
		Language: 1,
	}

	jsonBook := marshalCheckNoError(t, testBook)

	w := execRequest("POST", "/api/v1/books", bytes.NewReader(jsonBook))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	checkErrorBodyNotEmpty(w, t)
}

func TestPostBookBadRequestForeignKeyErr(t *testing.T) {
	testBook := models.Book{
		Title:    "Post foreign key test",
		Year:     1996,
		Pages:    200,
		Author:   1,
		Genre:    1,
		Language: 999,
	}

	jsonBook := marshalCheckNoError(t, testBook)

	w := execRequest("POST", "/api/v1/books", bytes.NewReader(jsonBook))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	checkErrorBodyNotEmpty(w, t)
}

// PUT /books
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
}

func TestPutBookNotFoundBigPathId(t *testing.T) {
	testBook := models.Book{
		Title:    "Put book test",
		Year:     1996,
		Pages:    593,
		Author:   1,
		Genre:    1,
		Language: 1,
	}

	jsonBook := marshalCheckNoError(t, testBook)

	w := execRequest("PUT", "/api/v1/books/9999", bytes.NewReader(jsonBook))
	assert.Equal(t, http.StatusNotFound, w.Code)

	checkErrorBodyNotEmpty(w, t)
}

func TestPutBookBadRequestForeignKeyErr(t *testing.T) {
	testBook := models.Book{
		Title:    "Put foreign key test",
		Year:     1996,
		Pages:    593,
		Author:   1,
		Genre:    1,
		Language: 999,
	}

	jsonBook := marshalCheckNoError(t, testBook)

	w := execRequest("PUT", "/api/v1/books/1", bytes.NewReader(jsonBook))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	checkErrorBodyNotEmpty(w, t)
}

func TestPutBookBadRequestStringId(t *testing.T) {
	testBook := models.Book{
		Title:    "Put book test",
		Year:     1996,
		Pages:    593,
		Author:   1,
		Genre:    1,
		Language: 1,
	}

	jsonBook := marshalCheckNoError(t, testBook)

	w := execRequest("PUT", "/api/v1/books/string", bytes.NewReader(jsonBook))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	checkErrorBodyNotEmpty(w, t)
}

func TestPutBookBadRequestMalformedJSON(t *testing.T) {
	jsonStr := []byte(`{"title":"JSON Test","year":1996,"pages":200,"author":"Should be number","genre":1}`)
	w := execRequest("PUT", "/api/v1/books/1", bytes.NewReader(jsonStr))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	checkErrorBodyNotEmpty(w, t)
}

func TestPutBookBadRequestBadJSON(t *testing.T) {
	testBook := models.Book{
		Title: "Put book test",
		Year:  1996,
		Pages: 593,
	}

	jsonBook := marshalCheckNoError(t, testBook)

	w := execRequest("PUT", "/api/v1/books/1", bytes.NewReader(jsonBook))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	checkErrorBodyNotEmpty(w, t)
}

// PATCH /books
func TestPatchBookSuccess(t *testing.T) {
	jsonStr := []byte(`{"title":"Patch book test", "pages":999}`)
	w := execRequest("PATCH", "/api/v1/books/1", bytes.NewReader(jsonStr))
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestPatchBookNotFoundBigPathId(t *testing.T) {
	jsonStr := []byte(`{"title":"Patch book test", "pages":999}`)
	w := execRequest("PATCH", "/api/v1/books/9999", bytes.NewReader(jsonStr))
	assert.Equal(t, http.StatusNotFound, w.Code)

	checkErrorBodyNotEmpty(w, t)
}

func TestPatchBookBadRequestStringPathId(t *testing.T) {
	jsonStr := []byte(`{"title":"Patch book test", "pages":999}`)
	w := execRequest("PATCH", "/api/v1/books/string", bytes.NewReader(jsonStr))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	checkErrorBodyNotEmpty(w, t)
}

func TestPatchBookBadRequestBadJSON(t *testing.T) {
	jsonStr := []byte(`{"title":"Patch book test", "pages":Should be number"}`)
	w := execRequest("PATCH", "/api/v1/books/1", bytes.NewReader(jsonStr))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	checkErrorBodyNotEmpty(w, t)
}
