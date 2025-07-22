package handler_test

import (
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

	runTestErrors(t, "GET", "genres", getTests)
}

// POST /genres
func TestPostGenre_Success(t *testing.T) {
	testGenre := models.Genre{
		Name: "Post test",
	}

	var rGenre models.Genre
	jsonGenre := marshalCheckNoError(t, testGenre)
	w := execAndCheck(t, "POST", "/api/v1/genres", jsonGenre, http.StatusCreated, &rGenre)
	defer database.DelGenre(rGenre.ID)

	expLoc := fmt.Sprintf("/api/v1/genres/%v", rGenre.ID)
	assert.Equal(t, expLoc, w.Result().Header.Get("Location"))

	assert.NotZero(t, rGenre.ID, "Auto generatated, non zero ID")
	assert.Equal(t, testGenre.Name, rGenre.Name)
}

func TestPostGenre_BadRequest_MalformedJSON(t *testing.T) {
	jsonBytes := []byte(`{"name":999}`)
	execAndCheckError(t, "POST", "/api/v1/genres", jsonBytes, http.StatusBadRequest)
}

func TestPostGenre_Error(t *testing.T) {
	postTests := map[string]ErrorTests{
		"BadRequest_ValidationErr_EmptyName": {
			body: marshalCheckNoError(t, models.Genre{
				Name: "",
			}),
			query:  "",
			status: http.StatusBadRequest,
		},
	}

	runTestErrors(t, "POST", "genres", postTests)
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
	putTests := map[string]ErrorTests{
		"NotFound_BigPathID": {
			body: marshalCheckNoError(t, models.Genre{
				Name: "Put test",
			}),
			query:  "/9999",
			status: http.StatusNotFound,
		},
		"BadRequest_StringID": {
			body: marshalCheckNoError(t, models.Genre{
				Name: "Put test",
			}),
			query:  "/string",
			status: http.StatusBadRequest,
		},
		"BadRequest_BadJSON": {
			body:   marshalCheckNoError(t, models.Genre{}),
			query:  "/1",
			status: http.StatusBadRequest,
		},
	}

	runTestErrors(t, "PUT", "genres", putTests)
}

// DELETE /genres/id
func TestDeleteGenre_Success(t *testing.T) {
	newID, err := database.InsertGenre(models.Genre{
		Name: "Delete tester",
	})
	assert.NoError(t, err)

	newGenreLoc := fmt.Sprintf("/api/v1/genres/%v", newID)
	execAndCheck(t, "DELETE", newGenreLoc, nil, http.StatusNoContent, nil)

	_, err = database.GetGenre(newID)
	assert.ErrorIs(t, err, db.ErrNotFound)
}

func TestDeleteGenre_Error(t *testing.T) {
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

	runTestErrors(t, "DELETE", "genres", deleteTests)
}

// OPTIONS /genres
// OPTIONS /genres/id
func TestOptionsGenres_Success(t *testing.T) {
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
			[]string{"GET", "PUT", "DELETE", "OPTIONS"},
		},
	}

	runTestOptionsSuccess(t, "genres", optionsTests)
}
