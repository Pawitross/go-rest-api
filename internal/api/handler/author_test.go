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
func TestListAuthors_Success(t *testing.T) {
	w := execRequest("GET", "/api/v1/authors", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var rAuthors []models.Author
	decodeJSONBodyCheckEmpty(w, t, &rAuthors)
}

func TestListAuthors_BadRequest_UnknownParam(t *testing.T) {
	w := execRequest("GET", "/api/v1/authors?foo=bar", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	checkErrorBodyNotEmpty(w, t)
}

// GET /authors/id
func TestGetAuthor_Success(t *testing.T) {
	w := execRequest("GET", "/api/v1/authors/1", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var rAuthor models.Author
	decodeJSONBodyCheckEmpty(w, t, &rAuthor)

	assert.NotEmpty(t, rAuthor, "Author in the response body should not be empty")
	assert.NotEmpty(t, rAuthor.FirstName, "FirstName should not be empty")
	assert.NotEmpty(t, rAuthor.LastName, "LastName should not be empty")
}

func TestGetAuthor_Error(t *testing.T) {
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
			w := execRequest("GET", "/api/v1/authors"+tt.query, nil)
			assert.Equal(t, tt.status, w.Code)

			checkErrorBodyNotEmpty(w, t)
		})
	}
}

// POST /authors
func TestPostAuthor_Success(t *testing.T) {
	testAuthor := models.Author{
		FirstName: "Post",
		LastName:  "test",
	}

	jsonAuthor := marshalCheckNoError(t, testAuthor)
	w := execRequest("POST", "/api/v1/authors", bytes.NewReader(jsonAuthor))
	assert.Equal(t, http.StatusCreated, w.Code)

	var rAuthor models.Author
	decodeJSONBodyCheckEmpty(w, t, &rAuthor)
	defer database.DelAuthor(rAuthor.Id)

	expLoc := fmt.Sprintf("/api/v1/authors/%v", rAuthor.Id)
	assert.Equal(t, expLoc, w.Result().Header.Get("Location"))

	assert.NotZero(t, rAuthor.Id, "Auto generatated, non zero ID")
	assert.Equal(t, testAuthor.FirstName, rAuthor.FirstName)
	assert.Equal(t, testAuthor.LastName, rAuthor.LastName)
}

func TestPostAuthor_BadRequest_MalformedJSON(t *testing.T) {
	jsonBytes := []byte(`{"first_name":999,"last_name":"test"}`)
	w := execRequest("POST", "/api/v1/authors", bytes.NewReader(jsonBytes))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	checkErrorBodyNotEmpty(w, t)
}

func TestPostAuthor_Error(t *testing.T) {
	postTests := map[string]struct {
		testAuthor models.Author
		status     int
	}{
		"BadRequest_ValidationErr": {
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
func TestPutAuthor_Success(t *testing.T) {
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

func TestPutAuthor_BadRequest_MalformedJSON(t *testing.T) {
	jsonBytes := []byte(`{"first_name":999,"last_name":"test"}`)
	w := execRequest("PUT", "/api/v1/authors/1", bytes.NewReader(jsonBytes))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	checkErrorBodyNotEmpty(w, t)
}

func TestPutAuthor_Error(t *testing.T) {
	putTests := map[string]struct {
		testAuthor models.Author
		query      string
		status     int
	}{
		"NotFound_BigPathId": {
			models.Author{
				FirstName: "Put",
				LastName:  "test",
			},
			"/9999",
			http.StatusNotFound,
		},
		"BadRequest_StringId": {
			models.Author{
				FirstName: "Put",
				LastName:  "test",
			},
			"/string",
			http.StatusBadRequest,
		},
		"BadRequest_BadJSON": {
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
func TestPatchAuthor_Success(t *testing.T) {
	jsonBytes := []byte(`{"first_name":"Patch test"}`)
	w := execRequest("PATCH", "/api/v1/authors/1", bytes.NewReader(jsonBytes))
	assert.Equal(t, http.StatusNoContent, w.Code)

	author, _ := database.GetAuthor(1)
	assert.Equal(t, "Patch test", author.FirstName)
}

func TestPatchAuthor_Error(t *testing.T) {
	patchTests := map[string]struct {
		jsonStr string
		query   string
		status  int
	}{
		"NotFound_BigPathId": {
			`{"first_name":"Patch test"}`,
			"/9999",
			http.StatusNotFound,
		},
		"BadRequest_StringPathId": {
			`{"first_name":"Patch test"}`,
			"/string",
			http.StatusBadRequest,
		},
		"BadRequest_BadJSON": {
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
func TestDeleteAuthor_Success(t *testing.T) {
	newId, err := database.InsertAuthor(models.Author{FirstName: "Delete", LastName: "tester"})
	assert.NoError(t, err)

	newAuthorLoc := fmt.Sprintf("/api/v1/authors/%v", newId)

	w := execRequest("DELETE", newAuthorLoc, nil)
	assert.Equal(t, http.StatusNoContent, w.Code)

	_, err = database.GetAuthor(newId)
	assert.ErrorIs(t, err, db.ErrNotFound)
}

func TestDeleteAuthor_Error(t *testing.T) {
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
			fullUrl := "/api/v1/authors" + tt.query

			w := execRequest("DELETE", fullUrl, nil)
			assert.Equal(t, tt.status, w.Code)

			checkErrorBodyNotEmpty(w, t)
		})
	}
}
