package routes_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"pawrest/internal/api/routes"
	"pawrest/internal/db/mock"
	"pawrest/internal/yamlconfig"
)

var authRouteTests = []struct {
	method   string
	endpoint string
	body     []byte
}{
	{"GET", "/api/v1/books", nil},
	{"GET", "/api/v1/books/1", nil},
	{"POST", "/api/v1/books", []byte(`{"title":"Route post test","year":1996,"pages":200,"author":1,"genre":1,"language":1}`)},
	{"PUT", "/api/v1/books/1", []byte(`{"title":"Route put test","year":2025,"pages":30,"author":1,"genre":1,"language":1}`)},
	{"PATCH", "/api/v1/books/1", []byte(`{"title":"Route patch test"}`)},
	{"DELETE", "/api/v1/books/2", nil},
	{"OPTIONS", "/api/v1/books", nil},
	{"OPTIONS", "/api/v1/books/1", nil},

	{"GET", "/api/v1/authors", nil},
	{"GET", "/api/v1/authors/1", nil},
	{"POST", "/api/v1/authors", []byte(`{"first_name":"Route post", "last_name":"test"}`)},
	{"PUT", "/api/v1/authors/1", []byte(`{"first_name":"Route put", "last_name":"test"}`)},
	{"PATCH", "/api/v1/authors/1", []byte(`{"first_name":"Route patch", "last_name":"test"}`)},
	{"DELETE", "/api/v1/authors/2", nil},
	{"OPTIONS", "/api/v1/authors", nil},
	{"OPTIONS", "/api/v1/authors/1", nil},

	{"GET", "/api/v1/genres", nil},
	{"GET", "/api/v1/genres/1", nil},
	{"POST", "/api/v1/genres", []byte(`{"name":"Route post"}`)},
	{"PUT", "/api/v1/genres/1", []byte(`{"name":"Route put"}`)},
	{"DELETE", "/api/v1/genres/2", nil},
	{"OPTIONS", "/api/v1/genres", nil},
	{"OPTIONS", "/api/v1/genres/1", nil},

	{"GET", "/api/v1/languages", nil},
	{"GET", "/api/v1/languages/1", nil},
	{"POST", "/api/v1/languages", []byte(`{"name":"Route post"}`)},
	{"PUT", "/api/v1/languages/1", []byte(`{"name":"Route put"}`)},
	{"DELETE", "/api/v1/languages/2", nil},
	{"OPTIONS", "/api/v1/languages", nil},
	{"OPTIONS", "/api/v1/languages/1", nil},
}

func setupTestRouter() *gin.Engine {
	mockdb := mock.NewMockDatabase()

	gin.SetMode(gin.TestMode)
	r := gin.New()

	cfg := &yamlconfig.Config{
		Secret: "secret-jwt-string",
	}

	routes.Router(r, mockdb, cfg)
	return r
}

func execRequest(r *gin.Engine, method, target string, body io.Reader, authHeader string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, target, body)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	r.ServeHTTP(w, req)
	return w
}

func getToken(t *testing.T, r *gin.Engine, adminToken bool) (string, bool) {
	jsonIn := []byte(nil)

	if adminToken {
		jsonIn = []byte(`{"return_admin_token":true}`)
	} else {
		jsonIn = []byte(`{"return_admin_token":false}`)
	}

	w := execRequest(r, "POST", "/api/v1/login", bytes.NewReader(jsonIn), "")

	var body map[string]any
	err := json.NewDecoder(w.Body).Decode(&body)
	assert.NoError(t, err)

	str, ok := body["token"].(string)
	return str, ok
}

func TestRoutes_NoAuth(t *testing.T) {
	router := setupTestRouter()

	noAuthRouteTests := []struct {
		method   string
		endpoint string
		body     []byte
	}{
		{"GET", "/swagger/index.html", nil},
	}

	for _, tt := range noAuthRouteTests {
		t.Run(tt.method+tt.endpoint, func(t *testing.T) {
			w := execRequest(router, tt.method, tt.endpoint, bytes.NewReader(tt.body), "")

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func TestRoutes_NoToken(t *testing.T) {
	router := setupTestRouter()

	for _, tt := range authRouteTests {
		t.Run(tt.method+tt.endpoint, func(t *testing.T) {
			w := execRequest(router, tt.method, tt.endpoint, bytes.NewReader(tt.body), "")

			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}

func TestRoutes_NonAdminToken(t *testing.T) {
	router := setupTestRouter()
	token, ok := getToken(t, router, false)
	if !ok {
		t.Fatalf("Failed to assert token type")
	}

	for _, tt := range authRouteTests {
		t.Run(tt.method+tt.endpoint, func(t *testing.T) {
			w := execRequest(router, tt.method, tt.endpoint, bytes.NewReader(tt.body), "Bearer "+token)

			switch tt.method {
			case "GET":
				assert.Equal(t, http.StatusOK, w.Code)
			case "OPTIONS":
				assert.Equal(t, http.StatusNoContent, w.Code)
			default:
				assert.Equal(t, http.StatusForbidden, w.Code)
			}
		})
	}
}

func TestRoutes_AdminToken(t *testing.T) {
	router := setupTestRouter()
	token, ok := getToken(t, router, true)
	if !ok {
		t.Fatalf("Failed to assert token type")
	}

	for _, tt := range authRouteTests {
		t.Run(tt.method+tt.endpoint, func(t *testing.T) {
			w := execRequest(router, tt.method, tt.endpoint, bytes.NewReader(tt.body), "Bearer "+token)

			switch tt.method {
			case "GET":
				assert.Equal(t, http.StatusOK, w.Code)
			case "POST":
				assert.Equal(t, http.StatusCreated, w.Code)
			case "PUT", "PATCH", "DELETE", "OPTIONS":
				assert.Equal(t, http.StatusNoContent, w.Code)
			}
		})
	}
}
