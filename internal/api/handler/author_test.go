package handler_test

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"pawrest/internal/db"
	"pawrest/internal/models"
)

// GET /authors
func TestListAuthors_Success(t *testing.T) {
	var rAuthors []models.Author
	execAndCheck(t, "GET", "/api/v1/authors", nil, http.StatusOK, &rAuthors)
}

func TestListAuthors_BadRequest_UnknownParam(t *testing.T) {
	execAndCheckError(t, "GET", "/api/v1/authors?foo=bar", nil, http.StatusBadRequest)
}

// GET /authors/id
func TestGetAuthor_Success(t *testing.T) {
	var rAuthor models.Author
	execAndCheck(t, "GET", "/api/v1/authors/1", nil, http.StatusOK, &rAuthor)

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
			fullUrl := "/api/v1/authors" + tt.query
			execAndCheckError(t, "GET", fullUrl, nil, tt.status)
		})
	}
}

// POST /authors
func TestPostAuthor_Success(t *testing.T) {
	testAuthor := models.Author{
		FirstName: "Post",
		LastName:  "test",
		BirthYear: 1970,
		DeathYear: models.I64Ptr(2043),
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
	assert.Equal(t, testAuthor.BirthYear, rAuthor.BirthYear)
	assert.Equal(t, testAuthor.DeathYear, rAuthor.DeathYear)
}

func TestPostAuthor_BadRequest_MalformedJSON(t *testing.T) {
	jsonBytes := []byte(`{"first_name":999,"last_name":"test"}`)
	execAndCheckError(t, "POST", "/api/v1/authors", jsonBytes, http.StatusBadRequest)
}

func TestPostAuthor_Error(t *testing.T) {
	postTests := map[string]struct {
		testAuthor models.Author
		status     int
	}{
		"BadRequest_ValidationErr_EmptyName": {
			models.Author{
				FirstName: "",
				LastName:  "test",
				BirthYear: 1970,
				DeathYear: models.I64Ptr(2043),
			},
			http.StatusBadRequest,
		},
		"BadRequest_ValidationErr_BirthGreaterThanDeath": {
			models.Author{
				FirstName: "",
				LastName:  "test",
				BirthYear: 2000,
				DeathYear: models.I64Ptr(1900),
			},
			http.StatusBadRequest,
		},
	}

	for name, tt := range postTests {
		t.Run(name, func(t *testing.T) {
			jsonAuthor := marshalCheckNoError(t, tt.testAuthor)
			execAndCheckError(t, "POST", "/api/v1/authors", jsonAuthor, tt.status)
		})
	}
}

// PUT /authors/id
func TestPutAuthor_Success(t *testing.T) {
	testAuthor := models.Author{
		FirstName: "Put",
		LastName:  "test",
		BirthYear: 2000,
		DeathYear: nil,
	}

	jsonAuthor := marshalCheckNoError(t, testAuthor)
	execAndCheck(t, "PUT", "/api/v1/authors/1", jsonAuthor, http.StatusNoContent, nil)

	author, _ := database.GetAuthor(1)
	assert.Equal(t, testAuthor.FirstName, author.FirstName)
	assert.Equal(t, testAuthor.LastName, author.LastName)
	assert.Equal(t, testAuthor.BirthYear, author.BirthYear)
	assert.Equal(t, testAuthor.DeathYear, author.DeathYear)
}

func TestPutAuthor_BadRequest_MalformedJSON(t *testing.T) {
	jsonBytes := []byte(`{"first_name":999,"last_name":"test"}`)
	execAndCheckError(t, "PUT", "/api/v1/authors/1", jsonBytes, http.StatusBadRequest)
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
				BirthYear: 1970,
				DeathYear: models.I64Ptr(2043),
			},
			"/9999",
			http.StatusNotFound,
		},
		"BadRequest_StringId": {
			models.Author{
				FirstName: "Put",
				LastName:  "test",
				BirthYear: 1970,
				DeathYear: models.I64Ptr(2043),
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
			execAndCheckError(t, "PUT", fullUrl, jsonAuthor, tt.status)
		})
	}
}

// PATCH /authors/id
func TestPatchAuthor_Success(t *testing.T) {
	jsonBytes := []byte(`{"first_name":"Patch test", "death_year": 2025}`)
	execAndCheck(t, "PATCH", "/api/v1/authors/1", jsonBytes, http.StatusNoContent, nil)

	author, _ := database.GetAuthor(1)
	assert.Equal(t, "Patch test", author.FirstName)
	assert.NotEmpty(t, author.LastName)
	assert.NotEmpty(t, author.BirthYear)
	assert.Equal(t, int64(2025), *author.DeathYear)
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
			execAndCheckError(t, "PATCH", fullUrl, []byte(tt.jsonStr), tt.status)
		})
	}
}

// DELETE /authors/id
func TestDeleteAuthor_Success(t *testing.T) {
	newId, err := database.InsertAuthor(models.Author{
		FirstName: "Delete",
		LastName:  "tester",
		BirthYear: 1900,
		DeathYear: nil,
	})
	assert.NoError(t, err)

	newAuthorLoc := fmt.Sprintf("/api/v1/authors/%v", newId)
	execAndCheck(t, "DELETE", newAuthorLoc, nil, http.StatusNoContent, nil)

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
			execAndCheckError(t, "DELETE", fullUrl, nil, tt.status)
		})
	}
}

// OPTIONS /authors
// OPTIONS /authors/id
func TestOptionsAuthors_Success(t *testing.T) {
	optionsTests := map[string]struct {
		query   string
		methods []string
	}{
		"BaseResource": {
			"",
			[]string{"GET", "POST", "OPTIONS"},
		},
		"PathId": {
			"/1",
			[]string{"GET", "PUT", "PATCH", "DELETE", "OPTIONS"},
		},
	}

	for name, tt := range optionsTests {
		t.Run(name, func(t *testing.T) {
			fullUrl := "/api/v1/authors" + tt.query

			w := execRequest("OPTIONS", fullUrl, nil)
			assert.Equal(t, http.StatusNoContent, w.Code)

			allowResp := w.Result().Header.Get("Allow")
			splitAllow := strings.Split(allowResp, ", ")

			assert.ElementsMatch(t, tt.methods, splitAllow)
		})
	}
}
