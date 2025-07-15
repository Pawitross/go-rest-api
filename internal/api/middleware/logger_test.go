package middleware_test

import (
	"bytes"
	"encoding/csv"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"pawrest/internal/api/middleware"
)

func fileExists(fName string) (bool, error) {
	if _, err := os.Stat(fName); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func removeLogFile(t *testing.T) {
	t.Helper()
	if exists, _ := fileExists("log.csv"); exists {
		if err := os.Remove("log.csv"); err != nil {
			t.Errorf("Failed to remove log file: %v", err)
		}
	}
}

func setupTestLoggingRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(middleware.FileLogger())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusCreated)
	})

	return router
}

func TestLogFileCreation(t *testing.T) {
	if err := middleware.InitLogger(); err != nil {
		t.Fatalf("Failed to initialize logging middleware: %v\n", err)
	}
	defer removeLogFile(t)
	defer middleware.CloseLogger()

	exists, err := fileExists("log.csv")
	if err != nil {
		t.Errorf("Error checking file existence: %v", err)
	}

	if !exists {
		t.Errorf("Log file doesn't exist, but should")
	}
}

func TestLogging(t *testing.T) {
	if err := middleware.InitLogger(); err != nil {
		t.Fatalf("Failed to initialize logging middleware: %v\n", err)
	}
	defer removeLogFile(t)
	defer middleware.CloseLogger()

	router := setupTestLoggingRouter()

	requests := []struct {
		method   string
		endpoint string
		query    string
		body     []byte
		status   int
	}{
		{"GET", "/test", "", nil, http.StatusCreated},
		{"POST", "/anothertest", "", nil, http.StatusNotFound},
		{"GET", "/querytest", "?foo=bar&baz=qux", nil, http.StatusNotFound},
		{"PATCH", "/testbody", "", []byte("Lorem ipsum"), http.StatusNotFound},
	}

	for _, r := range requests {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(r.method, r.endpoint+r.query, bytes.NewReader(r.body))
		req.RemoteAddr = "192.168.0.135:4444"
		router.ServeHTTP(w, req)

		assert.Equal(t, r.status, w.Code)
	}

	numberOfSentReq := len(requests)

	f, err := os.OpenFile("log.csv", os.O_RDONLY, 0644)
	if err != nil {
		t.Errorf("Error opening log file: %v", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ';'

	records, err := r.ReadAll()
	assert.NoError(t, err)

	assert.Equal(t, numberOfSentReq, len(records))

	for i, r := range requests {
		assert.NotEmpty(t, records[i][0])
		assert.Equal(t, "192.168.0.135", records[i][1])
		assert.Equal(t, r.method, records[i][2])
		assert.Equal(t, r.endpoint, records[i][3])
		assert.Equal(t, strings.TrimPrefix(r.query, "?"), records[i][4])
		assert.Equal(t, strconv.Itoa(r.status), records[i][5])
	}
}
