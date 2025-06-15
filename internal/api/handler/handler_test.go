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

func TestMain(m *testing.M) {
	if err := db.ConnectToDB(); err != nil {
		log.Fatal(err)
	}
	defer db.CloseDB()

	m.Run()
}

func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	routes.Router(router)

	return router
}

// GET /books
func TestListBooksSuccess(t *testing.T) {
	router := SetupTestRouter()
	w := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/api/v1/books", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var rBooks []models.Book

	err := json.NewDecoder(w.Body).Decode(&rBooks)
	assert.NoError(t, err, "Error decoding response data")

	assert.NotEmpty(t, rBooks, "Book slice should not be empty")
}

// GET /books/id
func TestGetBookSuccess(t *testing.T) {
	router := SetupTestRouter()
	w := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/api/v1/books/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var rBook models.Book
	err := json.NewDecoder(w.Body).Decode(&rBook)
	assert.NoError(t, err, "Error decoding response data")

	assert.NotEmpty(t, rBook, "Book in the response body should not be empty")
	assert.NotEmpty(t, rBook.Tytul, "Tytul should not be empty")
	assert.NotEmpty(t, rBook.Autor, "Autor should not be empty")
	assert.NotEmpty(t, rBook.Gatunek, "Gatunek should not be empty")
	assert.NotEmpty(t, rBook.Jezyk, "Jezyk should not be empty")
}

func TestGetBookNotFound(t *testing.T) {
	router := SetupTestRouter()
	w := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/api/v1/books/9999", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetBookBadRequest(t *testing.T) {
	router := SetupTestRouter()
	w := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/api/v1/books/string", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// POST /books
func TestPostBookSuccess(t *testing.T) {
	router := SetupTestRouter()
	w := httptest.NewRecorder()

	testBook := models.Book{
		Tytul:   "Post test book",
		Rok:     1996,
		Strony:  200,
		Autor:   1,
		Gatunek: 1,
		Jezyk:   1,
	}

	jsonBook, err := json.Marshal(testBook)
	assert.NoError(t, err, "Book JSON marshalling error")

	req := httptest.NewRequest("POST", "/api/v1/books", bytes.NewReader(jsonBook))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var rBook models.Book
	err = json.NewDecoder(w.Body).Decode(&rBook)
	assert.NoError(t, err, "Error decoding response data")
	defer db.DelBook(rBook.Id)

	expLoc := fmt.Sprintf("/api/v1/books/%v", rBook.Id)

	assert.Equal(t, expLoc, w.Result().Header.Get("Location"))

	assert.NotZero(t, rBook.Id, "Auto generatated, non zero ID")
	assert.Equal(t, testBook.Tytul, rBook.Tytul)
	assert.Equal(t, testBook.Rok, rBook.Rok)
	assert.Equal(t, testBook.Strony, rBook.Strony)
	assert.Equal(t, testBook.Autor, rBook.Autor)
	assert.Equal(t, testBook.Gatunek, rBook.Gatunek)
	assert.Equal(t, testBook.Jezyk, rBook.Jezyk)
}

func TestPostBookValidationErr(t *testing.T) {
	router := SetupTestRouter()
	w := httptest.NewRecorder()

	testBook := models.Book{
		Tytul:   "",
		Rok:     1996,
		Strony:  200,
		Autor:   1,
		Gatunek: 1,
		Jezyk:   1,
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
	router := SetupTestRouter()
	w := httptest.NewRecorder()

	testBook := models.Book{
		Tytul:   "Książka PUT",
		Rok:     1996,
		Strony:  593,
		Autor:   1,
		Gatunek: 1,
		Jezyk:   1,
	}

	jsonBook, err := json.Marshal(testBook)
	assert.NoError(t, err, "Book JSON marshalling error")

	req := httptest.NewRequest("PUT", "/api/v1/books/1", bytes.NewReader(jsonBook))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestPutBookNotFoundBigId(t *testing.T) {
	router := SetupTestRouter()
	w := httptest.NewRecorder()

	testBook := models.Book{
		Tytul:   "Książka PUT",
		Rok:     1996,
		Strony:  593,
		Autor:   1,
		Gatunek: 1,
		Jezyk:   1,
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
	router := SetupTestRouter()
	w := httptest.NewRecorder()

	testBook := models.Book{
		Tytul:   "Książka PUT",
		Rok:     1996,
		Strony:  593,
		Autor:   1,
		Gatunek: 1,
		Jezyk:   1,
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
	router := SetupTestRouter()
	w := httptest.NewRecorder()

	testBook := models.Book{
		Tytul:   "Książka PUT",
		Rok:     1996,
		Strony:  593,
		Gatunek: 1,
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
