package handler_test

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"pawrest/internal/db"
	"pawrest/internal/models"
)

// GET /books
func TestListBooks_Success(t *testing.T) {
	var rBooks []models.Book
	execAndCheck(t, "GET", "/api/v1/books", nil, http.StatusOK, &rBooks)
}

func TestListBooksExt_Success(t *testing.T) {
	var rBooks []models.BookExt
	execAndCheck(t, "GET", "/api/v1/books?extend=true", nil, http.StatusOK, &rBooks)
}

func TestListBooks_BadRequest_UnknownParam(t *testing.T) {
	execAndCheckError(t, "GET", "/api/v1/books?foo=bar", nil, http.StatusBadRequest)
}

// GET /books/id
func TestGetBook_Success(t *testing.T) {
	var rBook models.Book
	execAndCheck(t, "GET", "/api/v1/books/1", nil, http.StatusOK, &rBook)

	assert.NotEmpty(t, rBook, "Book in the response body should not be empty")
	assert.NotEmpty(t, rBook.Title, "Title should not be empty")
	assert.NotEmpty(t, rBook.Year, "Year should not be empty")
	assert.NotEmpty(t, rBook.Pages, "Pages should not be empty")
	assert.NotEmpty(t, rBook.Author, "Author should not be empty")
	assert.NotEmpty(t, rBook.Genre, "Genre should not be empty")
	assert.NotEmpty(t, rBook.Language, "Language should not be empty")
}

func TestGetBook_Error(t *testing.T) {
	getTests := map[string]struct {
		query  string
		status int
	}{
		"NotFound_BigPathId": {
			"/9999",
			http.StatusNotFound,
		},
		"BadRequest_StringPathId": {
			"/string",
			http.StatusBadRequest,
		},
	}

	for name, tt := range getTests {
		t.Run(name, func(t *testing.T) {
			fullUrl := "/api/v1/books" + tt.query
			execAndCheckError(t, "GET", fullUrl, nil, tt.status)
		})
	}
}

// POST /books
func TestPostBook_Success(t *testing.T) {
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
	decodeJSONBodyCheckEmpty(w, t, &rBook)
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

func TestPostBook_BadRequest_MalformedJSON(t *testing.T) {
	jsonBytes := []byte(`{"title":"JSON Test","year":1996,"pages":200,"author":"Should be number","genre":1}`)
	execAndCheckError(t, "POST", "/api/v1/books", jsonBytes, http.StatusBadRequest)
}

func TestPostBook_Error(t *testing.T) {
	postTests := map[string]struct {
		testBook models.Book
		status   int
	}{
		"BadRequest_ValidationErr": {
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
		"BadRequest_ForeignKeyErr": {
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
			execAndCheckError(t, "POST", "/api/v1/books", jsonBook, tt.status)
		})
	}
}

// PUT /books/id
func TestPutBook_Success(t *testing.T) {
	testBook := models.Book{
		Title:    "Put book test",
		Year:     1996,
		Pages:    593,
		Author:   1,
		Genre:    1,
		Language: 1,
	}

	jsonBook := marshalCheckNoError(t, testBook)
	execAndCheck(t, "PUT", "/api/v1/books/1", jsonBook, http.StatusNoContent, nil)

	book, _ := database.GetBook(1)
	assert.Equal(t, testBook.Title, book.Title)
	assert.Equal(t, testBook.Year, book.Year)
	assert.Equal(t, testBook.Pages, book.Pages)
	assert.Equal(t, testBook.Author, book.Author)
	assert.Equal(t, testBook.Genre, book.Genre)
	assert.Equal(t, testBook.Language, book.Language)
}

func TestPutBook_BadRequest_MalformedJSON(t *testing.T) {
	jsonBytes := []byte(`{"title":"JSON Test","year":1996,"pages":200,"author":"Should be number","genre":1}`)
	execAndCheckError(t, "PUT", "/api/v1/books/1", jsonBytes, http.StatusBadRequest)
}

func TestPutBook_Error(t *testing.T) {
	putTests := map[string]struct {
		testBook models.Book
		query    string
		status   int
	}{
		"NotFound_BigPathId": {
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
		"BadRequest_ForeignKeyErr": {
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
		"BadRequest_StringId": {
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
		"BadRequest_BadJSON": {
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
			execAndCheckError(t, "PUT", fullUrl, jsonBook, tt.status)
		})
	}
}

// PATCH /books/id
func TestPatchBook_Success(t *testing.T) {
	jsonBytes := []byte(`{"title":"Patch book test", "pages":999}`)
	execAndCheck(t, "PATCH", "/api/v1/books/1", jsonBytes, http.StatusNoContent, nil)

	book, _ := database.GetBook(1)
	assert.Equal(t, "Patch book test", book.Title)
	assert.Equal(t, int64(999), book.Pages)
}

func TestPatchBook_Error(t *testing.T) {
	patchTests := map[string]struct {
		jsonStr string
		query   string
		status  int
	}{
		"NotFound_BigPathId": {
			`{"title":"Patch book test", "pages":999}`,
			"/9999",
			http.StatusNotFound,
		},
		"BadRequest_StringPathId": {
			`{"title":"Patch book test", "pages":999}`,
			"/string",
			http.StatusBadRequest,
		},
		"BadRequest_BadJSON": {
			`{"title":"Patch book test", "pages":Should be number"}`,
			"/1",
			http.StatusBadRequest,
		},
	}

	for name, tt := range patchTests {
		t.Run(name, func(t *testing.T) {
			fullUrl := "/api/v1/books" + tt.query
			execAndCheckError(t, "PATCH", fullUrl, []byte(tt.jsonStr), tt.status)
		})
	}
}

// DELETE /books/id
func TestDeleteBook_Success(t *testing.T) {
	execAndCheck(t, "DELETE", "/api/v1/books/2", nil, http.StatusNoContent, nil)

	_, err := database.GetBook(2)
	assert.ErrorIs(t, err, db.ErrNotFound)
}

func TestDeleteBook_Error(t *testing.T) {
	deleteTests := map[string]struct {
		query  string
		status int
	}{
		"NotFound_BigPathId": {
			"/9999",
			http.StatusNotFound,
		},
		"BadRequest_StringPathId": {
			"/string",
			http.StatusBadRequest,
		},
	}

	for name, tt := range deleteTests {
		t.Run(name, func(t *testing.T) {
			fullUrl := "/api/v1/books" + tt.query
			execAndCheckError(t, "DELETE", fullUrl, nil, tt.status)
		})
	}
}
