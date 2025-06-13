package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"reflect"
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
	cfg.ClientFoundRows = true

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

func queryWithParams[T any](
	query string,
	params url.Values,
	allowPar map[string]string,
	scanFunc func(*T, *sql.Rows) error,
) ([]T, error) {
	var args []any

	if len(params) > 0 {
		filter, argsOut, err := assembleFilter(params, allowPar)
		if err != nil {
			return nil, err
		}

		query += filter
		args = argsOut
	}

	var data []T

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("Błąd zapytania (%v)", err)
	}
	defer rows.Close()

	for rows.Next() {
		var d T

		if err := scanFunc(&d, rows); err != nil {
			return nil, fmt.Errorf("Błąd odczytywania (%v)", err)
		}

		data = append(data, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Błąd (%v)", err)
	}

	return data, nil
}

func queryId[T any](
	query string,
	id int64,
	scanFunc func(*T, *sql.Row) error,
) (T, error) {
	var d T

	row := db.QueryRow(query, id)
	if err := scanFunc(&d, row); err != nil {
		if err == sql.ErrNoRows {
			return d, fmt.Errorf("Brak zasobu o id %d", id)
		}

		return d, fmt.Errorf("Błąd odczytywania (%v)", err)
	}

	return d, nil
}

func insert(query string, args ...any) (int64, error) {
	res, err := db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("Nie udało się dodać rekordu (%v)", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Nie udało się pobrać id (%v)", err)
	}

	return id, nil
}

func updateWholeId(query string, args ...any) error {
	res, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("Nie udało się zaktualizować (%v)", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Zmienione wiersze (%v)", err)
	}

	if rows == 0 {
		return fmt.Errorf("Nie znaleziono rekordu do aktualizacji")
	}

	return nil
}

func updatePartId(d any, table string, id int64, fToDb map[string]string) error {
	var (
		updates []string
		args    []any
	)

	valOfD := reflect.ValueOf(d)

	fields := reflect.VisibleFields(reflect.TypeOf(d))
	for i, f := range fields {
		fValue := valOfD.Field(i)
		if fValue.IsZero() {
			continue
		}

		columnName, ok := fToDb[f.Name]
		if !ok {
			continue
		}

		updates = append(updates, columnName+" = ?")
		args = append(args, fValue.Interface())
	}

	if len(updates) == 0 {
		return fmt.Errorf("Brak kolumn do zaktualizowania")
	}

	query := "UPDATE " + table + " SET " + strings.Join(updates, ", ") + " WHERE id = ?"
	args = append(args, id)

	res, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("Nie udało się zaktualizować (%v)", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Zmienione wiersze (%v)", err)
	}

	if rows == 0 {
		return fmt.Errorf("Nie znaleziono rekordu do aktualizacji")
	}

	return nil
}

func deleteId(query string, id int64) error {
	res, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Nie udało się usunąć (%v)", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Zmienione wiersze (%v)", err)
	}

	if rows == 0 {
		return fmt.Errorf("Brak zasobu o id %v", id)
	}

	return nil
}
