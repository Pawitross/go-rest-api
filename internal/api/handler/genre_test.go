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

// GET /genres
func TestListGenres_Success(t *testing.T) {
	var rGenres []models.Genre
	execAndCheck(t, "GET", "/api/v1/genres", nil, http.StatusOK, &rGenres)
}

func TestListGenres_BadRequest_UnknownParam(t *testing.T) {
	execAndCheckError(t, "GET", "/api/v1/genres?foo=bar", nil, http.StatusBadRequest)
}

// GET /genres/id
func TestGetGenre_Success(t *testing.T) {
	var rGenre models.Genre
	execAndCheck(t, "GET", "/api/v1/genres/1", nil, http.StatusOK, &rGenre)

	assert.NotEmpty(t, rGenre, "Genre in the response body should not be empty")
	assert.NotEmpty(t, rGenre.Name, "Name should not be empty")
}

func TestGetGenre_Error(t *testing.T) {
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
			fullUrl := "/api/v1/genres" + tt.query
			execAndCheckError(t, "GET", fullUrl, nil, tt.status)
		})
	}
}

// POST /genres
func TestPostGenre_Success(t *testing.T) {
	testGenre := models.Genre{
		Name: "Post test",
	}

	jsonGenre := marshalCheckNoError(t, testGenre)
	w := execRequest("POST", "/api/v1/genres", bytes.NewReader(jsonGenre))
	assert.Equal(t, http.StatusCreated, w.Code)

	var rGenre models.Genre
	decodeJSONBodyCheckEmpty(w, t, &rGenre)
	defer database.DelGenre(rGenre.Id)

	expLoc := fmt.Sprintf("/api/v1/genres/%v", rGenre.Id)
	assert.Equal(t, expLoc, w.Result().Header.Get("Location"))

	assert.NotZero(t, rGenre.Id, "Auto generatated, non zero ID")
	assert.Equal(t, testGenre.Name, rGenre.Name)
}

func TestPostGenre_BadRequest_MalformedJSON(t *testing.T) {
	jsonBytes := []byte(`{"name":999}`)
	execAndCheckError(t, "POST", "/api/v1/genres", jsonBytes, http.StatusBadRequest)
}

func TestPostGenre_Error(t *testing.T) {
	postTests := map[string]struct {
		testGenre models.Genre
		status    int
	}{
		"BadRequest_ValidationErr_EmptyName": {
			models.Genre{
				Name: "",
			},
			http.StatusBadRequest,
		},
	}

	for name, tt := range postTests {
		t.Run(name, func(t *testing.T) {
			jsonGenre := marshalCheckNoError(t, tt.testGenre)
			execAndCheckError(t, "POST", "/api/v1/genres", jsonGenre, tt.status)
		})
	}
}

// PUT /genres/id
func TestPutGenre_Success(t *testing.T) {
	testGenre := models.Genre{
		Name: "Put test",
	}

	jsonGenre := marshalCheckNoError(t, testGenre)
	execAndCheck(t, "PUT", "/api/v1/genres/1", jsonGenre, http.StatusNoContent, nil)

	genre, _ := database.GetGenre(1)
	assert.Equal(t, testGenre.Name, genre.Name)
}

func TestPutGenre_BadRequest_MalformedJSON(t *testing.T) {
	jsonBytes := []byte(`{"name":999}`)
	execAndCheckError(t, "PUT", "/api/v1/genres/1", jsonBytes, http.StatusBadRequest)
}

func TestPutGenre_Error(t *testing.T) {
	putTests := map[string]struct {
		testGenre models.Genre
		query     string
		status    int
	}{
		"NotFound_BigPathId": {
			models.Genre{
				Name: "Put test",
			},
			"/9999",
			http.StatusNotFound,
		},
		"BadRequest_StringId": {
			models.Genre{
				Name: "Put test",
			},
			"/string",
			http.StatusBadRequest,
		},
		"BadRequest_BadJSON": {
			models.Genre{},
			"/1",
			http.StatusBadRequest,
		},
	}

	for name, tt := range putTests {
		t.Run(name, func(t *testing.T) {
			fullUrl := "/api/v1/genres" + tt.query
			jsonGenre := marshalCheckNoError(t, tt.testGenre)
			execAndCheckError(t, "PUT", fullUrl, jsonGenre, tt.status)
		})
	}
}

// DELETE /genres/id
func TestDeleteGenre_Success(t *testing.T) {
	newId, err := database.InsertGenre(models.Genre{
		Name: "Delete tester",
	})
	assert.NoError(t, err)

	newGenreLoc := fmt.Sprintf("/api/v1/genres/%v", newId)
	execAndCheck(t, "DELETE", newGenreLoc, nil, http.StatusNoContent, nil)

	_, err = database.GetGenre(newId)
	assert.ErrorIs(t, err, db.ErrNotFound)
}

func TestDeleteGenre_Error(t *testing.T) {
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
			fullUrl := "/api/v1/genres" + tt.query
			execAndCheckError(t, "DELETE", fullUrl, nil, tt.status)
		})
	}
}
