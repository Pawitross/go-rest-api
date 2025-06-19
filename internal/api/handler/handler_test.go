package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"pawrest/internal/api/routes"
	"pawrest/internal/db"
	"pawrest/internal/models"
)

var database *db.Database

func TestMain(m *testing.M) {
	var err error
	database, err = db.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.CloseDB()

	m.Run()
}

func SetupTestRouter(db *db.Database) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	routes.Router(router, db)

	return router
}

// GET /books
func TestListBooksSuccess(t *testing.T) {
	router := SetupTestRouter(database)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/api/v1/books", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var rBooks []models.Book

	err := json.NewDecoder(w.Body).Decode(&rBooks)
	assert.NoError(t, err, "Error decoding response data")

	assert.NotEmpty(t, rBooks, "Book slice should not be empty")
}

func TestListBooksBadRequestUnknownParam(t *testing.T) {
	router := SetupTestRouter(database)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/api/v1/books?foo=bar", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var rError models.Error

	err := json.NewDecoder(w.Body).Decode(&rError)
	assert.NoError(t, err, "Error decoding response data")
}

// GET /books/id
func TestGetBookSuccess(t *testing.T) {
	router := SetupTestRouter(database)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/api/v1/books/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var rBook models.Book
	err := json.NewDecoder(w.Body).Decode(&rBook)
	assert.NoError(t, err, "Error decoding response data")

	assert.NotEmpty(t, rBook, "Book in the response body should not be empty")
	assert.NotEmpty(t, rBook.Title, "Title should not be empty")
	assert.NotEmpty(t, rBook.Author, "Author should not be empty")
	assert.NotEmpty(t, rBook.Genre, "Genre should not be empty")
	assert.NotEmpty(t, rBook.Language, "Language should not be empty")
}

func TestGetBookNotFound(t *testing.T) {
	router := SetupTestRouter(database)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/api/v1/books/9999", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetBookBadRequest(t *testing.T) {
	router := SetupTestRouter(database)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/api/v1/books/string", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// POST /books
func TestPostBookSuccess(t *testing.T) {
	router := SetupTestRouter(database)
	w := httptest.NewRecorder()

	testBook := models.Book{
		Title:    "Post test book",
		Year:     1996,
		Pages:    200,
		Author:   1,
		Genre:    1,
		Language: 1,
	}

	jsonBook, err := json.Marshal(testBook)
	assert.NoError(t, err, "Book JSON marshalling error")

	req := httptest.NewRequest("POST", "/api/v1/books", bytes.NewReader(jsonBook))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var rBook models.Book
	err = json.NewDecoder(w.Body).Decode(&rBook)
	assert.NoError(t, err, "Error decoding response data")
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

func TestPostBookValidationErr(t *testing.T) {
	router := SetupTestRouter(database)
	w := httptest.NewRecorder()

	testBook := models.Book{
		Title:    "",
		Year:     1996,
		Pages:    200,
		Author:   1,
		Genre:    1,
		Language: 1,
	}

	jsonBook, err := json.Marshal(testBook)
	assert.NoError(t, err, "Book JSON marshalling error")

	req := httptest.NewRequest("POST", "/api/v1/books", bytes.NewReader(jsonBook))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var rError models.Error
	err = json.NewDecoder(w.Body).Decode(&rError)
	assert.NoError(t, err, "Error decoding response data")

	assert.NotEmpty(t, rError.Error, "Error message should not be empty")
}

// PUT /books
func TestPutBookSuccess(t *testing.T) {
	router := SetupTestRouter(database)
	w := httptest.NewRecorder()

	testBook := models.Book{
		Title:    "Książka PUT",
		Year:     1996,
		Pages:    593,
		Author:   1,
		Genre:    1,
		Language: 1,
	}

	jsonBook, err := json.Marshal(testBook)
	assert.NoError(t, err, "Book JSON marshalling error")

	req := httptest.NewRequest("PUT", "/api/v1/books/1", bytes.NewReader(jsonBook))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestPutBookNotFoundBigId(t *testing.T) {
	router := SetupTestRouter(database)
	w := httptest.NewRecorder()

	testBook := models.Book{
		Title:    "Książka PUT",
		Year:     1996,
		Pages:    593,
		Author:   1,
		Genre:    1,
		Language: 1,
	}

	jsonBook, err := json.Marshal(testBook)
	assert.NoError(t, err, "Book JSON marshalling error")

	req := httptest.NewRequest("PUT", "/api/v1/books/9999", bytes.NewReader(jsonBook))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var rError models.Error
	err = json.NewDecoder(w.Body).Decode(&rError)
	assert.NoError(t, err, "Error decoding response data")

	assert.NotEmpty(t, rError.Error, "Error message should not be empty")
}

func TestPutBookBadRequestStringId(t *testing.T) {
	router := SetupTestRouter(database)
	w := httptest.NewRecorder()

	testBook := models.Book{
		Title:    "Książka PUT",
		Year:     1996,
		Pages:    593,
		Author:   1,
		Genre:    1,
		Language: 1,
	}

	jsonBook, err := json.Marshal(testBook)
	assert.NoError(t, err, "Book JSON marshalling error")

	req := httptest.NewRequest("PUT", "/api/v1/books/string", bytes.NewReader(jsonBook))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var rError models.Error
	err = json.NewDecoder(w.Body).Decode(&rError)
	assert.NoError(t, err, "Error decoding response data")

	assert.NotEmpty(t, rError.Error, "Error message should not be empty")
}

func TestPutBookBadRequestBadJSON(t *testing.T) {
	router := SetupTestRouter(database)
	w := httptest.NewRecorder()

	testBook := models.Book{
		Title: "Książka PUT",
		Year:  1996,
		Pages: 593,
		Genre: 1,
	}

	jsonBook, err := json.Marshal(testBook)
	assert.NoError(t, err, "Book JSON marshalling error")

	req := httptest.NewRequest("PUT", "/api/v1/books/1", bytes.NewReader(jsonBook))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var rError models.Error
	err = json.NewDecoder(w.Body).Decode(&rError)
	assert.NoError(t, err, "Error decoding response data")

	assert.NotEmpty(t, rError.Error, "Error message should not be empty")
}
