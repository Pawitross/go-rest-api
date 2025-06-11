package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	m "pawrest/internal/models"
)

func GetBooks(params url.Values) ([]m.Book, error) {
	query := "SELECT id, tytul, rok_wydania, liczba_stron, id_autora, id_gatunku, id_jezyka FROM ksiazka"

	allowedParams := map[string]string{
		"id":       "id",
		"title":    "tytul",
		"year":     "rok_wydania",
		"pages":    "liczba_stron",
		"author":   "id_autora",
		"genre":    "id_gatunku",
		"language": "id_jezyka",
	}

	bookFunc := func(b *m.Book, rows *sql.Rows) error {
		return rows.Scan(&b.Id, &b.Tytul, &b.Rok, &b.Strony, &b.Autor, &b.Gatunek, &b.Jezyk)
	}

	return queryWithParams[m.Book](
		query,
		params,
		allowedParams,
		bookFunc,
	)
}

func bookExists(id int64) error {
	// TODO
	return fmt.Errorf("TODO")
}

func GetBook(id int64) (m.Book, error) {
	query := "SELECT id, tytul, rok_wydania, liczba_stron, id_autora, id_gatunku, id_jezyka FROM ksiazka WHERE id = ?"

	bookFunc := func(b *m.Book, row *sql.Row) error {
		return row.Scan(&b.Id, &b.Tytul, &b.Rok, &b.Strony, &b.Autor, &b.Gatunek, &b.Jezyk)
	}

	return queryId[m.Book](query, id, bookFunc)
}

func InsertBook(b m.Book) (int64, error) {
	query := "INSERT INTO ksiazka (tytul, rok_wydania, liczba_stron, id_autora, id_gatunku, id_jezyka) VALUES (?, ?, ?, ?, ?, ?)"

	return insert(query, b.Tytul, b.Rok, b.Strony, b.Autor, b.Gatunek, b.Jezyk)
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

	var (
		updates []string
		args    []any
	)

	valOfB := reflect.ValueOf(b)

	fields := reflect.VisibleFields(reflect.TypeOf(b))
	for i, f := range fields {
		fValue := valOfB.Field(i)
		if fValue.IsZero() {
			continue
		}

		columnName, ok := fieldToDb[f.Name]
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
