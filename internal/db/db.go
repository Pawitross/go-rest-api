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

var Db *sql.DB

func ConnectToDB() error {
	cfg := mysql.NewConfig()

	cfg.User = "root"
	cfg.Passwd = ""
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "paw"

	var err error
	Db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}

	if pingErr := Db.Ping(); pingErr != nil {
		return pingErr
	}

	//fmt.Printf("Połączono z bazą danych %q pod adresem %q\n", cfg.DBName, cfg.Addr)
	return nil
}

func CloseDB() {
	if Db != nil {
		Db.Close()
	}
}

func assembleFilter(params url.Values, allowedParams map[string]string) (string, []any, error) {
	filter := ""
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

	limit, limitted := params["limit"]
	if limitted {
		delete(params, "limit")
	}

	offset, offsetted := params["offset"]
	if offsetted {
		delete(params, "offset")
	}

	if !limitted && offsetted {
		return "", nil, fmt.Errorf("Nie podano limitu do podanego offsetu.")
	}

	for k, vsl := range params {
		op := "="

		for kop, vop := range operators {
			if before, queryop := strings.CutSuffix(k, kop); queryop {
				k = before
				op = vop
				break
			}
		}

		dbCol, allowed := allowedParams[k]
		if !allowed || vsl[0] == "" || len(vsl) == 0 {
			return "", nil, fmt.Errorf("Wprowadzono nieznany parametr lub jest pusty.")
		}

		if len(vsl) > 1 {
			return "", nil, fmt.Errorf("Wprowadzono za dużo parametrów dla jednej kolumny.")
		}

		conditions = append(conditions, dbCol+" "+op+" ?")
		args = append(args, vsl[0])
	}

	if len(conditions) > 0 {
		filter += " WHERE " + strings.Join(conditions, " AND ")
	}

	if limitted {
		filter += " LIMIT ?"
		args = append(args, limit[0])
	}

	if offsetted {
		filter += " OFFSET ?"
		args = append(args, offset[0])
	}

	//fmt.Println("conditions:", conditions)
	//fmt.Println("filter:", filter)
	//fmt.Println("args:", args)
	//fmt.Println()

	return filter, args, nil
}

func GetKsiazki(params url.Values) ([]m.Ksiazka, error) {
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

	var ksiazki []m.Ksiazka

	rows, err := Db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("Błąd zapytania (%v)", err)
	}
	defer rows.Close()

	for rows.Next() {
		var k m.Ksiazka

		if err := rows.Scan(&k.Id, &k.Tytul, &k.Rok, &k.Strony, &k.Autor, &k.Gatunek, &k.Jezyk); err != nil {
			return nil, fmt.Errorf("Błąd odczytywania (%v)", err)
		}

		ksiazki = append(ksiazki, k)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Błąd (%v)", err)
	}

	return ksiazki, nil
}

func ksiazkaExists(id int64) error {
	// TODO
	return fmt.Errorf("TODO")
}

func GetKsiazka(id int64) (m.Ksiazka, error) {
	var k m.Ksiazka

	query := "SELECT id, tytul, rok_wydania, liczba_stron, id_autora, id_gatunku, id_jezyka FROM ksiazka WHERE id = ?"

	row := Db.QueryRow(query, id)
	if err := row.Scan(&k.Id, &k.Tytul, &k.Rok, &k.Strony, &k.Autor, &k.Gatunek, &k.Jezyk); err != nil {
		if err == sql.ErrNoRows {
			return k, fmt.Errorf("Brak książki o id %d", id)
		}

		return k, fmt.Errorf("Błąd odczytywania (%v)", err)
	}

	return k, nil
}

func InsertKsiazka(k m.Ksiazka) (int64, error) {
	query := "INSERT INTO ksiazka (tytul, rok_wydania, liczba_stron, id_autora, id_gatunku, id_jezyka) VALUES (?, ?, ?, ?, ?, ?)"

	res, err := Db.Exec(query, k.Tytul, k.Rok, k.Strony, k.Autor, k.Gatunek, k.Jezyk)
	if err != nil {
		return 0, fmt.Errorf("Nie udało się dodać rekordu (%v)", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Nie udało się pobrać id (%v)", err)
	}

	return id, nil
}

func UpdateWholeKsiazka(id int64, k m.Ksiazka) error {
	query := "UPDATE ksiazka SET tytul = ?, rok_wydania = ?, liczba_stron = ?, id_autora = ?, id_gatunku = ?, id_jezyka = ? WHERE id = ?"

	res, err := Db.Exec(query, k.Tytul, k.Rok, k.Strony, k.Autor, k.Gatunek, k.Jezyk, id)
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

func UpdateKsiazka(id int64, k m.Ksiazka) error {
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

	valOfK := reflect.ValueOf(k)

	fields := reflect.VisibleFields(reflect.TypeOf(k))
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

	res, err := Db.Exec(query, args...)
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

func DelKsiazka(id int64) error {
	query := "DELETE FROM ksiazka WHERE id = ?"

	res, err := Db.Exec(query, id)
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
