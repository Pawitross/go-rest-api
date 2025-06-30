package middleware_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"pawrest/internal/api/handler"
	"pawrest/internal/api/middleware"
	"pawrest/internal/models"
)

func setupTestAuthRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/authenticate", middleware.Authenticate(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "You're in!"})
	})

	router.GET("/authorize", middleware.Authenticate(), middleware.Authorize(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "You're in!"})
	})

	router.POST("/login", handler.ReturnToken)

	return router
}

func getToken(t *testing.T, r *gin.Engine, adminToken bool) string {
	jsonIn := []byte(nil)

	if adminToken {
		jsonIn = []byte(`{"return_admin_token":true}`)
	} else {
		jsonIn = []byte(`{"return_admin_token":false}`)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(jsonIn))
	r.ServeHTTP(w, req)

	var body models.Token
	err := json.NewDecoder(w.Body).Decode(&body)
	assert.NoError(t, err)

	return body.Token
}

func TestAuthentication(t *testing.T) {
	router := setupTestAuthRouter()

	token := getToken(t, router, false)

	authTests := map[string]struct {
		header string
		status int
		errMsg string
	}{
		"NoHeader": {
			"",
			http.StatusUnauthorized,
			"No token provided",
		},
		"FooHeader": {
			"foo",
			http.StatusUnauthorized,
			"No Bearer prefix in Authorization header",
		},
		"OnlyBearerPrefixHeader": {
			"Bearer ",
			http.StatusUnauthorized,
			"Malformed token",
		},
		"InvalidToken": {
			"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30",
			http.StatusUnauthorized,
			"Invalid token signature",
		},
		"ExpiredToken": {
			"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsImV4cCI6MTc0OTQ3MTM1NSwiaWF0IjoxNzQ5NDY5NTU1LCJpc3MiOiJzZXJ2ZXIiLCJzdWIiOiJ1c2VyIn0.WSGMrO3uGYIYNdvLZzbk2x4K0KZFmFS0H4PmyeP8kHY",
			http.StatusUnauthorized,
			"Token has expired",
		},
		"ValidToken": {
			"Bearer " + token,
			http.StatusOK,
			"",
		},
	}

	for name, tt := range authTests {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/authenticate", nil)

			if tt.header != "" {
				req.Header.Add("Authorization", tt.header)
			}

			router.ServeHTTP(w, req)
			assert.Equal(t, tt.status, w.Code)

			if tt.errMsg != "" {
				var rError models.Error
				err := json.NewDecoder(w.Body).Decode(&rError)
				assert.NoError(t, err)

				assert.NotEmpty(t, rError)
				assert.Equal(t, tt.errMsg, rError.Error)
			} else {
				var rResp map[string]string
				err := json.NewDecoder(w.Body).Decode(&rResp)
				assert.NoError(t, err)

				assert.NotEmpty(t, rResp)
				assert.Equal(t, "You're in!", rResp["message"])
			}
		})
	}
}

func TestAuthorization(t *testing.T) {
	router := setupTestAuthRouter()

	userToken := getToken(t, router, false)
	adminToken := getToken(t, router, true)

	authTests := map[string]struct {
		token  string
		status int
		errMsg string
	}{
		"NonAdminToken": {
			userToken,
			http.StatusForbidden,
			"You do not have sufficient permissions to access this resource",
		},
		"AdminToken": {
			adminToken,
			http.StatusOK,
			"",
		},
	}

	for name, tt := range authTests {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/authorize", nil)
			req.Header.Add("Authorization", "Bearer "+tt.token)

			router.ServeHTTP(w, req)
			assert.Equal(t, tt.status, w.Code)

			if tt.errMsg != "" {
				var rError models.Error
				err := json.NewDecoder(w.Body).Decode(&rError)
				assert.NoError(t, err)

				assert.NotEmpty(t, rError)
				assert.Equal(t, tt.errMsg, rError.Error)
			} else {
				var rResp map[string]string
				err := json.NewDecoder(w.Body).Decode(&rResp)
				assert.NoError(t, err)

				assert.NotEmpty(t, rResp)
				assert.Equal(t, "You're in!", rResp["message"])
			}
		})
	}
}
