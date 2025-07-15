package handler_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"pawrest/internal/db"
	"pawrest/internal/models"
)

// GET /languages
func TestListLanguages_Success(t *testing.T) {
	var rLanguages []models.Language
	execAndCheck(t, "GET", "/api/v1/languages", nil, http.StatusOK, &rLanguages)
}

func TestListLanguages_BadRequest_UnknownParam(t *testing.T) {
	execAndCheckError(t, "GET", "/api/v1/languages?foo=bar", nil, http.StatusBadRequest)
}

// GET /languages/id
func TestGetLanguage_Success(t *testing.T) {
	var rLanguage models.Language
	execAndCheck(t, "GET", "/api/v1/languages/1", nil, http.StatusOK, &rLanguage)

	assert.NotEmpty(t, rLanguage, "Language in the response body should not be empty")
	assert.NotEmpty(t, rLanguage.Name, "Name should not be empty")
}

func TestGetLanguage_Error(t *testing.T) {
	getTests := map[string]ErrorTests{
		"NotFound_BigPathId": {
			body:   nil,
			query:  "/9999",
			status: http.StatusNotFound,
		},
		"BadRequest_StringPathId": {
			body:   nil,
			query:  "/string",
			status: http.StatusBadRequest,
		},
	}

	runTestErrors(t, "GET", "languages", getTests)
}

// POST /languages
func TestPostLanguage_Success(t *testing.T) {
	testLanguage := models.Language{
		Name: "Post test",
	}

	var rLanguage models.Language
	jsonLanguage := marshalCheckNoError(t, testLanguage)
	w := execAndCheck(t, "POST", "/api/v1/languages", jsonLanguage, http.StatusCreated, &rLanguage)
	defer database.DelLanguage(rLanguage.Id)

	expLoc := fmt.Sprintf("/api/v1/languages/%v", rLanguage.Id)
	assert.Equal(t, expLoc, w.Result().Header.Get("Location"))

	assert.NotZero(t, rLanguage.Id, "Auto generatated, non zero ID")
	assert.Equal(t, testLanguage.Name, rLanguage.Name)
}

func TestPostLanguage_BadRequest_MalformedJSON(t *testing.T) {
	jsonBytes := []byte(`{"name":999}`)
	execAndCheckError(t, "POST", "/api/v1/languages", jsonBytes, http.StatusBadRequest)
}

func TestPostLanguage_Error(t *testing.T) {
	postTests := map[string]ErrorTests{
		"BadRequest_ValidationErr_EmptyName": {
			body: marshalCheckNoError(t, models.Language{
				Name: "",
			}),
			query:  "",
			status: http.StatusBadRequest,
		},
	}

	runTestErrors(t, "POST", "languages", postTests)
}

// PUT /languages/id
func TestPutLanguage_Success(t *testing.T) {
	testLanguage := models.Language{
		Name: "Put test",
	}

	jsonLanguage := marshalCheckNoError(t, testLanguage)
	execAndCheck(t, "PUT", "/api/v1/languages/1", jsonLanguage, http.StatusNoContent, nil)

	language, _ := database.GetLanguage(1)
	assert.Equal(t, testLanguage.Name, language.Name)
}

func TestPutLanguage_BadRequest_MalformedJSON(t *testing.T) {
	jsonBytes := []byte(`{"name":999}`)
	execAndCheckError(t, "PUT", "/api/v1/languages/1", jsonBytes, http.StatusBadRequest)
}

func TestPutLanguage_Error(t *testing.T) {
	putTests := map[string]ErrorTests{
		"NotFound_BigPathId": {
			body: marshalCheckNoError(t, models.Language{
				Name: "Put test",
			}),
			query:  "/9999",
			status: http.StatusNotFound,
		},
		"BadRequest_StringId": {
			body: marshalCheckNoError(t, models.Language{
				Name: "Put test",
			}),
			query:  "/string",
			status: http.StatusBadRequest,
		},
		"BadRequest_BadJSON": {
			body:   marshalCheckNoError(t, models.Language{}),
			query:  "/1",
			status: http.StatusBadRequest,
		},
	}

	runTestErrors(t, "PUT", "languages", putTests)
}

// DELETE /languages/id
func TestDeleteLanguage_Success(t *testing.T) {
	newId, err := database.InsertLanguage(models.Language{
		Name: "Delete tester",
	})
	assert.NoError(t, err)

	newLanguageLoc := fmt.Sprintf("/api/v1/languages/%v", newId)
	execAndCheck(t, "DELETE", newLanguageLoc, nil, http.StatusNoContent, nil)

	_, err = database.GetLanguage(newId)
	assert.ErrorIs(t, err, db.ErrNotFound)
}

func TestDeleteLanguage_Error(t *testing.T) {
	deleteTests := map[string]ErrorTests{
		"NotFound_BigPathId": {
			body:   nil,
			query:  "/9999",
			status: http.StatusNotFound,
		},
		"BadRequest_StringPathId": {
			body:   nil,
			query:  "/string",
			status: http.StatusBadRequest,
		},
	}

	runTestErrors(t, "DELETE", "languages", deleteTests)
}

// OPTIONS /languages
// OPTIONS /languages/id
func TestOptionsLanguages_Success(t *testing.T) {
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
			[]string{"GET", "PUT", "DELETE", "OPTIONS"},
		},
	}

	runTestOptionsSuccess(t, "languages", optionsTests)
}
