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

// GET /authors
func TestListAuthorsSuccess(t *testing.T) {
	w := execRequest("GET", "/api/v1/authors", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var rAuthors []models.Author
	decodeBodyCheckEmpty(w, t, &rAuthors)
}

func TestListAuthorsBadRequestUnknownParam(t *testing.T) {
	w := execRequest("GET", "/api/v1/authors?foo=bar", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	checkErrorBodyNotEmpty(w, t)
}

// GET /authors/id
func TestGetAuthorSuccess(t *testing.T) {
	w := execRequest("GET", "/api/v1/authors/1", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var rAuthor models.Author
	decodeBodyCheckEmpty(w, t, &rAuthor)

	assert.NotEmpty(t, rAuthor, "Author in the response body should not be empty")
	assert.NotEmpty(t, rAuthor.FirstName, "FirstName should not be empty")
	assert.NotEmpty(t, rAuthor.LastName, "LastName should not be empty")
}

func TestGetAuthorNotFoundBigPathId(t *testing.T) {
	w := execRequest("GET", "/api/v1/authors/9999", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)

	checkErrorBodyNotEmpty(w, t)
}

func TestGetAuthorBadRequestStringPathId(t *testing.T) {
	w := execRequest("GET", "/api/v1/authors/string", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	checkErrorBodyNotEmpty(w, t)
}

// POST /authors
func TestPostAuthorSuccess(t *testing.T) {
	testAuthor := models.Author{
		FirstName: "Post",
		LastName:  "test",
	}

	jsonAuthor := marshalCheckNoError(t, testAuthor)

	w := execRequest("POST", "/api/v1/authors", bytes.NewReader(jsonAuthor))
	assert.Equal(t, http.StatusCreated, w.Code)

	var rAuthor models.Author
	decodeBodyCheckEmpty(w, t, &rAuthor)
	defer database.DelAuthor(rAuthor.Id)

	expLoc := fmt.Sprintf("/api/v1/authors/%v", rAuthor.Id)
	assert.Equal(t, expLoc, w.Result().Header.Get("Location"))

	assert.NotZero(t, rAuthor.Id, "Auto generatated, non zero ID")
	assert.Equal(t, testAuthor.FirstName, rAuthor.FirstName)
	assert.Equal(t, testAuthor.LastName, rAuthor.LastName)
}

func TestPostAuthorBadRequestMalformedJSON(t *testing.T) {
	jsonStr := []byte(`{"first_name":999,"last_name":"test"}`)
	w := execRequest("POST", "/api/v1/authors", bytes.NewReader(jsonStr))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	checkErrorBodyNotEmpty(w, t)
}

func TestPostAuthorError(t *testing.T) {
	postTests := map[string]struct {
		testAuthor models.Author
		status     int
	}{
		"BadRequestValidationErr": {
			models.Author{
				FirstName: "",
				LastName:  "test",
			},
			http.StatusBadRequest,
		},
	}

	for name, tt := range postTests {
		t.Run(name, func(t *testing.T) {
			jsonAuthor := marshalCheckNoError(t, tt.testAuthor)

			w := execRequest("POST", "/api/v1/authors", bytes.NewReader(jsonAuthor))
			assert.Equal(t, tt.status, w.Code)

			checkErrorBodyNotEmpty(w, t)
		})
	}
}

// PUT /authors/id
func TestPutAuthorSuccess(t *testing.T) {
	testAuthor := models.Author{
		FirstName: "Put",
		LastName:  "test",
	}

	jsonAuthor := marshalCheckNoError(t, testAuthor)

	w := execRequest("PUT", "/api/v1/authors/1", bytes.NewReader(jsonAuthor))
	assert.Equal(t, http.StatusNoContent, w.Code)

	author, _ := database.GetAuthor(1)
	assert.Equal(t, testAuthor.FirstName, author.FirstName)
	assert.Equal(t, testAuthor.LastName, author.LastName)
}

func TestPutAuthorBadRequestMalformedJSON(t *testing.T) {
	jsonStr := []byte(`{"first_name":999,"last_name":"test"}`)
	w := execRequest("PUT", "/api/v1/authors/1", bytes.NewReader(jsonStr))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	checkErrorBodyNotEmpty(w, t)
}

func TestPutAuthorError(t *testing.T) {
	putTests := map[string]struct {
		testAuthor models.Author
		query      string
		status     int
	}{
		"NotFoundBigPathId": {
			models.Author{
				FirstName: "Put",
				LastName:  "test",
			},
			"/9999",
			http.StatusNotFound,
		},
		"BadRequestStringId": {
			models.Author{
				FirstName: "Put",
				LastName:  "test",
			},
			"/string",
			http.StatusBadRequest,
		},
		"BadRequestBadJSON": {
			models.Author{
				FirstName: "Put",
			},
			"/1",
			http.StatusBadRequest,
		},
	}

	for name, tt := range putTests {
		t.Run(name, func(t *testing.T) {
			fullUrl := "/api/v1/authors" + tt.query

			jsonAuthor := marshalCheckNoError(t, tt.testAuthor)

			w := execRequest("PUT", fullUrl, bytes.NewReader(jsonAuthor))
			assert.Equal(t, tt.status, w.Code)

			checkErrorBodyNotEmpty(w, t)
		})
	}
}

// PATCH /authors/id
func TestPatchAuthorSuccess(t *testing.T) {
	jsonStr := []byte(`{"first_name":"Patch test"}`)
	w := execRequest("PATCH", "/api/v1/authors/1", bytes.NewReader(jsonStr))
	assert.Equal(t, http.StatusNoContent, w.Code)

	author, _ := database.GetAuthor(1)
	assert.Equal(t, "Patch test", author.FirstName)
}

func TestPatchAuthorError(t *testing.T) {
	patchTests := map[string]struct {
		jsonStr string
		query   string
		status  int
	}{
		"NotFoundBigPathId": {
			`{"first_name":"Patch test"}`,
			"/9999",
			http.StatusNotFound,
		},
		"BadRequestStringPathId": {
			`{"first_name":"Patch test"}`,
			"/string",
			http.StatusBadRequest,
		},
		"BadRequestBadJSON": {
			`{"first_name":9999999}`,
			"/1",
			http.StatusBadRequest,
		},
	}

	for name, tt := range patchTests {
		t.Run(name, func(t *testing.T) {
			fullUrl := "/api/v1/authors" + tt.query

			w := execRequest("PATCH", fullUrl, bytes.NewReader([]byte(tt.jsonStr)))
			assert.Equal(t, tt.status, w.Code)

			checkErrorBodyNotEmpty(w, t)
		})
	}
}

// DELETE /authors/id
func TestDeleteAuthorSuccess(t *testing.T) {
	newId, err := database.InsertAuthor(models.Author{FirstName: "Delete", LastName: "tester"})
	assert.NoError(t, err)

	newAuthorLoc := fmt.Sprintf("/api/v1/authors/%v", newId)

	w := execRequest("DELETE", newAuthorLoc, nil)
	assert.Equal(t, http.StatusNoContent, w.Code)

	_, err = database.GetAuthor(newId)
	assert.ErrorIs(t, err, db.ErrNotFound)
}

func TestDeleteAuthorError(t *testing.T) {
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
			fullUrl := "/api/v1/authors" + tt.query

			w := execRequest("DELETE", fullUrl, nil)
			assert.Equal(t, tt.status, w.Code)

			checkErrorBodyNotEmpty(w, t)
		})
	}
}
