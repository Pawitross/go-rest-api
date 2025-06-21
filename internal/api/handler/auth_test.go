package handler_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"pawrest/internal/models"
)

func checkTokenStructure(t *testing.T, token string) {
	assert.NotEmpty(t, token)

	assert.True(t, strings.HasPrefix(token, "eyJ"))
	assert.Equal(t, 2, strings.Count(token, "."))
}

func TestLoginTokenSuccess(t *testing.T) {
	urltests := []struct {
		query string
	}{
		{""},
		{"?admin=false"},
		{"?admin=true"},
	}

	for _, tt := range urltests {
		t.Run(tt.query, func(t *testing.T) {
			fullUrl := "/api/v1/login" + tt.query

			w := execRequest("GET", fullUrl, nil)
			assert.Equal(t, http.StatusOK, w.Code)

			var rToken models.Token
			decodeBodyCheckEmpty(w, t, &rToken)

			token := rToken.Token
			checkTokenStructure(t, token)
		})
	}
}

func TestLoginTokenBadRequest(t *testing.T) {
	w := execRequest("GET", "/api/v1/login?admin=foo", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	checkErrorBodyNotEmpty(w, t)
}
