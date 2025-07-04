package handler_test

import (
	"bytes"
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

func TestLoginToken_Success(t *testing.T) {
	loginTests := []string{
		`{"return_admin_token":false}`,
		`{"return_admin_token":true}`,
	}

	for _, tc := range loginTests {
		t.Run(tc, func(t *testing.T) {
			w := execRequest("POST", "/api/v1/login", bytes.NewReader([]byte(tc)))
			assert.Equal(t, http.StatusOK, w.Code)

			var rToken models.Token
			decodeJSONBodyCheckEmpty(w, t, &rToken)

			token := rToken.Token
			checkTokenStructure(t, token)
		})
	}
}

func TestLoginToken_BadRequest(t *testing.T) {
	errorTests := []string{
		`{}`,
		`{"admin":}`,
		`{"admin":true}`,
		`{"admin":false}`,
		`{"return_admin_token":FALSE}`,
		`{"return_admin_token":TRUE}`,
		`{"return_admin_token":0}`,
		`{"return_admin_token":1}`,
		`{"return_admin_token":falseFoo}`,
		`{"return_admin_token":trueFoo}`,
		`{"return_admin_token":}`,
		`{"return_admin_token":""}`,
		`{"return_admin_token":"false"}`,
		`{"return_admin_token":"true"}`,
	}

	for _, tc := range errorTests {
		t.Run(tc, func(t *testing.T) {
			w := execRequest("POST", "/api/v1/login", bytes.NewReader([]byte(tc)))
			assert.Equal(t, http.StatusBadRequest, w.Code)

			checkErrorBodyNotEmpty(w, t)
		})
	}
}
