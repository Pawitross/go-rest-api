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
			fullUrl := "/api/v1/languages" + tt.query
			execAndCheckError(t, "GET", fullUrl, nil, tt.status)
		})
	}
}

// POST /languages
func TestPostLanguage_Success(t *testing.T) {
	testLanguage := models.Language{
		Name: "Post test",
	}

	jsonLanguage := marshalCheckNoError(t, testLanguage)
	w := execRequest("POST", "/api/v1/languages", bytes.NewReader(jsonLanguage))
	assert.Equal(t, http.StatusCreated, w.Code)

	var rLanguage models.Language
	decodeJSONBodyCheckEmpty(w, t, &rLanguage)
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
	postTests := map[string]struct {
		testLanguage models.Language
		status       int
	}{
		"BadRequest_ValidationErr_EmptyName": {
			models.Language{
				Name: "",
			},
			http.StatusBadRequest,
		},
	}

	for name, tt := range postTests {
		t.Run(name, func(t *testing.T) {
			jsonLanguage := marshalCheckNoError(t, tt.testLanguage)
			execAndCheckError(t, "POST", "/api/v1/languages", jsonLanguage, tt.status)
		})
	}
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
	putTests := map[string]struct {
		testLanguage models.Language
		query        string
		status       int
	}{
		"NotFound_BigPathId": {
			models.Language{
				Name: "Put test",
			},
			"/9999",
			http.StatusNotFound,
		},
		"BadRequest_StringId": {
			models.Language{
				Name: "Put test",
			},
			"/string",
			http.StatusBadRequest,
		},
		"BadRequest_BadJSON": {
			models.Language{},
			"/1",
			http.StatusBadRequest,
		},
	}

	for name, tt := range putTests {
		t.Run(name, func(t *testing.T) {
			fullUrl := "/api/v1/languages" + tt.query
			jsonLanguage := marshalCheckNoError(t, tt.testLanguage)
			execAndCheckError(t, "PUT", fullUrl, jsonLanguage, tt.status)
		})
	}
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
			fullUrl := "/api/v1/languages" + tt.query
			execAndCheckError(t, "DELETE", fullUrl, nil, tt.status)
		})
	}
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

	for name, tt := range optionsTests {
		t.Run(name, func(t *testing.T) {
			fullUrl := "/api/v1/languages" + tt.query

			w := execRequest("OPTIONS", fullUrl, nil)
			assert.Equal(t, http.StatusNoContent, w.Code)

			allowResp := w.Result().Header.Get("Allow")
			splitAllow := strings.Split(allowResp, ", ")

			assert.ElementsMatch(t, tt.methods, splitAllow)
		})
	}
}
