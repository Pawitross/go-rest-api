package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/go-sql-driver/mysql"
	m "pawrest/internal/models"
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
	conditions := []string{}
	args := []any{}
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

func GetBooks(params url.Values) ([]m.Book, error) {
	query := "SELECT id, tytul, rok_wydania, liczba_stron, id_autora, id_gatunku, id_jezyka FROM ksiazka"
	args := []any{}

	if len(params) > 0 {
		allowedParams := map[string]string{
			"id":       "id",
			"title":    "tytul",
			"year":     "rok_wydania",
			"pages":    "liczba_stron",
			"author":   "id_autora",
			"genre":    "id_gatunku",
			"language": "id_jezyka",
		}

		filter, argsOut, err := assembleFilter(params, allowedParams)
		if err != nil {
			return nil, err
		}

		query += filter
		args = argsOut
	}

	var books []m.Book

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("Błąd zapytania (%v)", err)
	}
	defer rows.Close()

	for rows.Next() {
		var b m.Book

		if err := rows.Scan(&b.Id, &b.Tytul, &b.Rok, &b.Strony, &b.Autor, &b.Gatunek, &b.Jezyk); err != nil {
			return nil, fmt.Errorf("Błąd odczytywania (%v)", err)
		}

		books = append(books, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Błąd (%v)", err)
	}

	return books, nil
}

func bookExists(id int64) error {
	// TODO
	return fmt.Errorf("TODO")
}

func GetBook(id int64) (m.Book, error) {
	var b m.Book

	query := "SELECT id, tytul, rok_wydania, liczba_stron, id_autora, id_gatunku, id_jezyka FROM ksiazka WHERE id = ?"

	row := db.QueryRow(query, id)
	if err := row.Scan(&b.Id, &b.Tytul, &b.Rok, &b.Strony, &b.Autor, &b.Gatunek, &b.Jezyk); err != nil {
		if err == sql.ErrNoRows {
			return b, fmt.Errorf("Brak książki o id %d", id)
		}

		return b, fmt.Errorf("Błąd odczytywania (%v)", err)
	}

	return b, nil
}

func InsertBook(b m.Book) (int64, error) {
	query := "INSERT INTO ksiazka (tytul, rok_wydania, liczba_stron, id_autora, id_gatunku, id_jezyka) VALUES (?, ?, ?, ?, ?, ?)"

	res, err := db.Exec(query, b.Tytul, b.Rok, b.Strony, b.Autor, b.Gatunek, b.Jezyk)
	if err != nil {
		return 0, fmt.Errorf("Nie udało się dodać rekordu (%v)", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Nie udało się pobrać id (%v)", err)
	}

	return id, nil
}

func UpdateWholeBook(id int64, b m.Book) error {
	query := "UPDATE ksiazka SET tytul = ?, rok_wydania = ?, liczba_stron = ?, id_autora = ?, id_gatunku = ?, id_jezyka = ? WHERE id = ?"

	res, err := db.Exec(query, b.Tytul, b.Rok, b.Strony, b.Autor, b.Gatunek, b.Jezyk, id)
	if err != nil {
		return fmt.Errorf("Nie udało się zaktualizować (%v)", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Zmienione wiersze (%v)", err)
	}

	if rows == 0 {
		return fmt.Errorf("Nie znaleziono rekordu do aktualizacji lub nie zmieniono rekordu")
	}

	return nil
}

func UpdateBook(id int64, b m.Book) error {
	fieldToDb := map[string]string{
		"Tytul":   "tytul",
		"Rok":     "rok_wydania",
		"Strony":  "liczba_stron",
		"Autor":   "id_autora",
		"Gatunek": "id_gatunku",
		"Jezyk":   "id_jezyka",
	}

	updates := []string{}
	args := []any{}

	valOfK := reflect.ValueOf(b)

	fields := reflect.VisibleFields(reflect.TypeOf(b))
	for i, field := range fields {
		fValue := valOfK.Field(i)
		if fValue.IsZero() {
			continue
		}

		columnName, ok := fieldToDb[field.Name]
		if !ok {
			continue
		}

		updates = append(updates, columnName+" = ?")
		args = append(args, fValue.Interface())
	}

	if len(updates) == 0 {
		return fmt.Errorf("Brak kolumn do zaktualizowania")
	}

	query := "UPDATE ksiazka SET " + strings.Join(updates, ", ") + " WHERE id = ?"
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
		return fmt.Errorf("Nie znaleziono rekordu do aktualizacji lub nie zmieniono rekordu")
	}

	return nil
}

func DelBook(id int64) error {
	query := "DELETE FROM ksiazka WHERE id = ?"

	res, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Nie udało się usunąć (%v)", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Zmienione wiersze (%v)", err)
	}

	if rows == 0 {
		return fmt.Errorf("Brak ksiązki o id %v", id)
	}

	return nil
}
