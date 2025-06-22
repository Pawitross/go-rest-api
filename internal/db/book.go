package db

import (
	"database/sql"
	"fmt"
	"net/url"

	"pawrest/internal/models"
)

type BookDatabaseInterface interface {
	GetBooks(params url.Values) ([]models.Book, error)
	GetBooksExt(params url.Values) ([]models.BookExt, error)
	GetBook(id int64) (models.Book, error)
	InsertBook(b models.Book) (int64, error)
	UpdateWholeBook(id int64, b models.Book) error
	UpdateBook(id int64, b models.Book) error
	DelBook(id int64) error
}

func (d *Database) GetBooks(params url.Values) ([]models.Book, error) {
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

	bookFunc := func(b *models.Book, rows *sql.Rows) error {
		return rows.Scan(&b.Id, &b.Title, &b.Year, &b.Pages, &b.Author, &b.Genre, &b.Language)
	}

	return queryWithParams[models.Book](
		d,
		query,
		params,
		allowedParams,
		bookFunc,
	)
}

func (d *Database) GetBooksExt(params url.Values) ([]models.BookExt, error) {
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

	bookFunc := func(b *models.BookExt, rows *sql.Rows) error {
		return rows.Scan(
			&b.Id,
			&b.Title,
			&b.Year,
			&b.Pages,
			&b.Author.Id,
			&b.Author.FirstName,
			&b.Author.LastName,
			&b.Genre.Id,
			&b.Genre.Name,
			&b.Language.Id,
			&b.Language.Name,
		)
	}

	return queryWithParams[models.BookExt](
		d,
		query,
		params,
		allowedParams,
		bookFunc,
	)
}

func (d *Database) bookExists(id int64) error {
	// TODO
	return fmt.Errorf("TODO")
}

func (d *Database) GetBook(id int64) (models.Book, error) {
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

	bookFunc := func(b *models.Book, row *sql.Row) error {
		return row.Scan(&b.Id, &b.Title, &b.Year, &b.Pages, &b.Author, &b.Genre, &b.Language)
	}

	return queryId[models.Book](d, query, id, bookFunc)
}

func (d *Database) InsertBook(b models.Book) (int64, error) {
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

	return d.insert(query, b.Title, b.Year, b.Pages, b.Author, b.Genre, b.Language)
}

func (d *Database) UpdateWholeBook(id int64, b models.Book) error {
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

	return d.updateWholeId(query, b.Title, b.Year, b.Pages, b.Author, b.Genre, b.Language, id)
}

func (d *Database) UpdateBook(id int64, b models.Book) error {
	fieldToDb := map[string]string{
		"Title":    "tytul",
		"Year":     "rok_wydania",
		"Pages":    "liczba_stron",
		"Author":   "id_autora",
		"Genre":    "id_gatunku",
		"Language": "id_jezyka",
	}

	return d.updatePartId(b, "ksiazka", id, fieldToDb)
}

func (d *Database) DelBook(id int64) error {
	query := "DELETE FROM ksiazka WHERE id = ?"

	return d.deleteId(query, id)
}
