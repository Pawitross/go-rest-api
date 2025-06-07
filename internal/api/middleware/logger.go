package middleware

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const timeFormat = "02.01.2006 15:04:05"

var logFile *os.File

type logRecord struct {
	time    string
	ip_addr string
	method  string
	path    string
	query   string
	status  int
}

func InitLogger() error {
	fileName := "log.csv"

	var err error
	logFile, err = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	return err
}

func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}

func logToFile(r logRecord) error {
	strRecord := []string{
		r.time,
		r.ip_addr,
		r.method,
		r.path,
		r.query,
		strconv.Itoa(r.status),
	}

	w := csv.NewWriter(logFile)
	w.Comma = ';'

	if err := w.Write(strRecord); err != nil {
		return err
	}

	w.Flush()
	return nil
}

func FileLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		record := logRecord{
			time.Now().Format(timeFormat),
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.URL.RawQuery,
			c.Writer.Status(),
		}

		if err := logToFile(record); err != nil {
			log.Println(err)
		}
	}
}
