package db

import (
	"database/sql"
	"fmt"
	"net/url"

	m "pawrest/internal/models"
)

func GetBooks(params url.Values) ([]m.Book, error) {
	query := `
	SELECT
		id,
		tytul,
		rok_wydania,
		liczba_stron,
		id_autora,
		id_gatunku,
		id_jezyka
	FROM ksiazka`

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

func GetBooksExt(params url.Values) ([]m.BookExt, error) {
	query := `
	SELECT
		k.id,
		tytul,
		rok_wydania,
		liczba_stron,
		id_autora,
		a.imie,
		a.nazwisko,
		id_gatunku,
		g.nazwa,
		id_jezyka,
		j.nazwa
	FROM ksiazka k
		JOIN autor a ON k.id_autora = a.id
		JOIN gatunek g ON k.id_gatunku = g.id
		JOIN jezyk j ON k.id_jezyka = j.id`

	allowedParams := map[string]string{
		"id":                "k.id",
		"title":             "tytul",
		"year":              "rok_wydania",
		"pages":             "liczba_stron",
		"author.id":         "id_autora",
		"author.first_name": "a.imie",
		"author.last_name":  "a.nazwisko",
		"genre.id":          "id_gatunku",
		"genre.name":        "g.nazwa",
		"language.id":       "id_jezyka",
		"language.name":     "j.nazwa",
	}

	bookFunc := func(b *m.BookExt, rows *sql.Rows) error {
		return rows.Scan(
			&b.Id,
			&b.Tytul,
			&b.Rok,
			&b.Strony,
			&b.Autor.Id,
			&b.Autor.Imie,
			&b.Autor.Nazwisko,
			&b.Gatunek.Id,
			&b.Gatunek.Nazwa,
			&b.Jezyk.Id,
			&b.Jezyk.Nazwa,
		)
	}

	return queryWithParams[m.BookExt](
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
	query := `
	SELECT
		id,
		tytul,
		rok_wydania,
		liczba_stron,
		id_autora,
		id_gatunku,
		id_jezyka
	FROM ksiazka
	WHERE id = ?`

	bookFunc := func(b *m.Book, row *sql.Row) error {
		return row.Scan(&b.Id, &b.Tytul, &b.Rok, &b.Strony, &b.Autor, &b.Gatunek, &b.Jezyk)
	}

	return queryId[m.Book](query, id, bookFunc)
}

func InsertBook(b m.Book) (int64, error) {
	query := `
	INSERT INTO ksiazka (
		tytul,
		rok_wydania,
		liczba_stron,
		id_autora,
		id_gatunku,
		id_jezyka
	)
	VALUES (?, ?, ?, ?, ?, ?)`

	return insert(query, b.Tytul, b.Rok, b.Strony, b.Autor, b.Gatunek, b.Jezyk)
}

func UpdateWholeBook(id int64, b m.Book) error {
	query := `
	UPDATE ksiazka
	SET
		tytul = ?,
		rok_wydania = ?,
		liczba_stron = ?,
		id_autora = ?,
		id_gatunku = ?,
		id_jezyka = ?
	WHERE id = ?`

	return updateWholeId(query, b.Tytul, b.Rok, b.Strony, b.Autor, b.Gatunek, b.Jezyk, id)
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

	return updatePartId(b, "ksiazka", id, fieldToDb)
}

func DelBook(id int64) error {
	query := "DELETE FROM ksiazka WHERE id = ?"

	return deleteId(query, id)
}
