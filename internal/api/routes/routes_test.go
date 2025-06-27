package routes_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"pawrest/internal/api/routes"
	"pawrest/internal/db/mock"
)

func setupTestRouter() *gin.Engine {
	mockdb := &mock.MockDatabase{}

	gin.SetMode(gin.TestMode)
	r := gin.New()

	routes.Router(r, mockdb)
	return r
}

func TestRoutes(t *testing.T) {
	router := setupTestRouter()

	routeTests := []struct {
		method   string
		endpoint string
		body     []byte
	}{
		{"GET", "/api/v1/books", nil},
		{"POST", "/api/v1/books", []byte(`{"title":"Route post test","year":1996,"pages":200,"author":1,"genre":1,"language":1}`)},
		{"GET", "/api/v1/books/1", nil},
		{"PUT", "/api/v1/books/1", []byte(`{"title":"Route put test","year":2025,"pages":30,"author":1,"genre":1,"language":1}`)},
		{"PATCH", "/api/v1/books/1", []byte(`{"title":"Route patch test"}`)},
		{"DELETE", "/api/v1/books/2", nil},
		{"GET", "/swagger/index.html", nil},
	}

	for _, tt := range routeTests {
		t.Run(tt.method+tt.endpoint, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.endpoint, bytes.NewReader(tt.body))

			router.ServeHTTP(w, req)

			assert.NotEqual(t, http.StatusNotFound, w.Code)
			assert.NotEqual(t, http.StatusBadRequest, w.Code)
			assert.NotEqual(t, http.StatusInternalServerError, w.Code)
		})
	}
}
