package middleware_test

import (
	"bytes"
	"encoding/csv"
	"errors"
	"log"
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

func runMain(m *testing.M) (int, error) {
	defer func() {
		if err := os.Remove("log.csv"); err != nil {
			log.Println(err)
		}
	}()

	return m.Run(), nil
}

func TestMain(m *testing.M) {
	code, err := runMain(m)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

func fileExists(fName string) (bool, error) {
	if _, err := os.Stat(fName); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, err
		}
		return true, err
	}
	return true, nil
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(middleware.FileLogger())
	router.GET("/test", func(c *gin.Context) {
		c.Status(201)
	})

	return router
}

func TestLogFileCreation(t *testing.T) {
	if err := middleware.InitLogger(); err != nil {
		t.Errorf("Failed to initialize logging middleware: %v\n", err)
	}
	defer middleware.CloseLogger()

	if exists, err := fileExists("log.csv"); !exists && err != nil {
		t.Errorf("File doesn't exist")
	}
}

func TestLogging(t *testing.T) {
	if err := middleware.InitLogger(); err != nil {
		t.Errorf("Failed to initialize logging middleware: %v\n", err)
	}
	defer middleware.CloseLogger()

	router := setupTestRouter()

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
