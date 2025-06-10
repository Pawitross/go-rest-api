package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func ConnectToDB() error {
	cfg := mysql.NewConfig()

	cfg.User = "root"
	cfg.Passwd = ""
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "paw"

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}

	if pingErr := db.Ping(); pingErr != nil {
		return pingErr
	}

	//fmt.Printf("Połączono z bazą danych %q pod adresem %q\n", cfg.DBName, cfg.Addr)
	return nil
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}

func assembleFilter(params url.Values, allowedParams map[string]string) (string, []any, error) {
	var (
		conditions []string
		args       []any
	)

	operators := map[string]string{
		".eq":  "=",
		".gt":  ">",
		".lt":  "<",
		".gte": ">=",
		".lte": "<=",
		".neq": "<>",
	}

	limit, hasLimit := params["limit"]
	offset, hasOffset := params["offset"]

	if !hasLimit && hasOffset {
		return "", nil, fmt.Errorf("Nie podano limitu do podanego offsetu.")
	}

	for key, valSlice := range params {
		if key == "limit" || key == "offset" {
			continue
		}

		operator := "="

		for suffix, sqlOp := range operators {
			if before, found := strings.CutSuffix(key, suffix); found {
				key = before
				operator = sqlOp
				break
			}
		}

		columnName, allowed := allowedParams[key]
		if !allowed {
			return "", nil, fmt.Errorf("Wprowadzono nieznany parametr.")
		}

		if len(valSlice) == 0 || valSlice[0] == "" {
			return "", nil, fmt.Errorf("Wprowadzony parametr jest pusty.")
		}

		if len(valSlice) > 1 {
			return "", nil, fmt.Errorf("Wprowadzono za dużo parametrów dla jednej kolumny.")
		}

		conditions = append(conditions, columnName+" "+operator+" ?")
		args = append(args, valSlice[0])
	}

	filter := ""

	if len(conditions) > 0 {
		filter += " WHERE " + strings.Join(conditions, " AND ")
	}

	if hasLimit {
		filter += " LIMIT ?"
		args = append(args, limit[0])
	}

	if hasOffset {
		filter += " OFFSET ?"
		args = append(args, offset[0])
	}

	return filter, args, nil
}
