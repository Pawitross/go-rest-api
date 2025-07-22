package handler_test

import (
	"fmt"
	"net/http"
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
	getTests := map[string]ErrorTests{
		"NotFound_BigPathID": {
			body:   nil,
			query:  "/9999",
			status: http.StatusNotFound,
		},
		"BadRequest_StringPathID": {
			body:   nil,
			query:  "/string",
			status: http.StatusBadRequest,
		},
	}

	runTestErrors(t, "GET", "authors", getTests)
}

// POST /authors
func TestPostAuthor_Success(t *testing.T) {
	testAuthor := models.Author{
		FirstName: "Post",
		LastName:  "test",
		BirthYear: 1970,
		DeathYear: models.I64Ptr(2043),
	}

	var rAuthor models.Author
	jsonAuthor := marshalCheckNoError(t, testAuthor)
	w := execAndCheck(t, "POST", "/api/v1/authors", jsonAuthor, http.StatusCreated, &rAuthor)
	defer database.DelAuthor(rAuthor.ID)

	expLoc := fmt.Sprintf("/api/v1/authors/%v", rAuthor.ID)
	assert.Equal(t, expLoc, w.Result().Header.Get("Location"))

	assert.NotZero(t, rAuthor.ID, "Auto generatated, non zero ID")
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
	postTests := map[string]ErrorTests{
		"BadRequest_ValidationErr_EmptyName": {
			body: marshalCheckNoError(t, models.Author{
				FirstName: "",
				LastName:  "test",
				BirthYear: 1970,
				DeathYear: models.I64Ptr(2043),
			}),
			query:  "",
			status: http.StatusBadRequest,
		},
		"BadRequest_ValidationErr_BirthGreaterThanDeath": {
			body: marshalCheckNoError(t, models.Author{
				FirstName: "",
				LastName:  "test",
				BirthYear: 2000,
				DeathYear: models.I64Ptr(1900),
			}),
			query:  "",
			status: http.StatusBadRequest,
		},
	}

	runTestErrors(t, "POST", "authors", postTests)
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
	putTests := map[string]ErrorTests{
		"NotFound_BigPathID": {
			body: marshalCheckNoError(t, models.Author{
				FirstName: "Put",
				LastName:  "test",
				BirthYear: 1970,
				DeathYear: models.I64Ptr(2043),
			}),
			query:  "/9999",
			status: http.StatusNotFound,
		},
		"BadRequest_StringID": {
			body: marshalCheckNoError(t, models.Author{
				FirstName: "Put",
				LastName:  "test",
				BirthYear: 1970,
				DeathYear: models.I64Ptr(2043),
			}),
			query:  "/string",
			status: http.StatusBadRequest,
		},
		"BadRequest_BadJSON": {
			body: marshalCheckNoError(t, models.Author{
				FirstName: "Put",
			}),
			query:  "/1",
			status: http.StatusBadRequest,
		},
	}

	runTestErrors(t, "PUT", "authors", putTests)
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
	patchTests := map[string]ErrorTests{
		"NotFound_BigPathID": {
			body:   []byte(`{"first_name":"Patch test"}`),
			query:  "/9999",
			status: http.StatusNotFound,
		},
		"BadRequest_StringPathID": {
			body:   []byte(`{"first_name":"Patch test"}`),
			query:  "/string",
			status: http.StatusBadRequest,
		},
		"BadRequest_BadJSON": {
			body:   []byte(`{"first_name":9999999}`),
			query:  "/1",
			status: http.StatusBadRequest,
		},
	}

	runTestErrors(t, "PATCH", "authors", patchTests)
}

// DELETE /authors/id
func TestDeleteAuthor_Success(t *testing.T) {
	newID, err := database.InsertAuthor(models.Author{
		FirstName: "Delete",
		LastName:  "tester",
		BirthYear: 1900,
		DeathYear: nil,
	})
	assert.NoError(t, err)

	newAuthorLoc := fmt.Sprintf("/api/v1/authors/%v", newID)
	execAndCheck(t, "DELETE", newAuthorLoc, nil, http.StatusNoContent, nil)

	_, err = database.GetAuthor(newID)
	assert.ErrorIs(t, err, db.ErrNotFound)
}

func TestDeleteAuthor_Error(t *testing.T) {
	deleteTests := map[string]ErrorTests{
		"NotFound_BigPathID": {
			body:   nil,
			query:  "/9999",
			status: http.StatusNotFound,
		},
		"BadRequest_StringPathID": {
			body:   nil,
			query:  "/string",
			status: http.StatusBadRequest,
		},
	}

	runTestErrors(t, "DELETE", "authors", deleteTests)
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
		"PathID": {
			"/1",
			[]string{"GET", "PUT", "PATCH", "DELETE", "OPTIONS"},
		},
	}

	runTestOptionsSuccess(t, "authors", optionsTests)
}
