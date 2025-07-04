package db

import (
	"database/sql"
	"net/url"

	"pawrest/internal/models"
)

type AuthorDatabaseInterface interface {
	GetAuthors(params url.Values) ([]models.Author, error)
	GetAuthor(id int64) (models.Author, error)
	InsertAuthor(a models.Author) (int64, error)
	UpdateWholeAuthor(id int64, a models.Author) error
	UpdateAuthor(id int64, a models.Author) error
	DelAuthor(id int64) error
}

func (d *Database) GetAuthors(params url.Values) ([]models.Author, error) {
	query := `
	SELECT id, imie, nazwisko
	FROM autor`

	allowedParams := map[string]string{
		"id":         "id",
		"first_name": "imie",
		"last_name":  "nazwisko",
	}

	authorFunc := func(a *models.Author, rows *sql.Rows) error {
		return rows.Scan(&a.Id, &a.FirstName, &a.LastName)
	}

	return queryWithParams[models.Author](
		d,
		query,
		params,
		allowedParams,
		authorFunc,
	)
}

func (d *Database) GetAuthor(id int64) (models.Author, error) {
	query := `
	SELECT id, imie, nazwisko
	FROM autor
	WHERE id = ?`

	authorFunc := func(a *models.Author, row *sql.Row) error {
		return row.Scan(&a.Id, &a.FirstName, &a.LastName)
	}

	return queryId[models.Author](d, query, id, authorFunc)
}

func (d *Database) InsertAuthor(a models.Author) (int64, error) {
	query := `
	INSERT INTO autor (imie, nazwisko)
	VALUES (?, ?)`

	return d.insert(query, a.FirstName, a.LastName)
}

func (d *Database) UpdateWholeAuthor(id int64, a models.Author) error {
	query := `
	UPDATE autor
	SET
		imie = ?,
		nazwisko = ?
	WHERE id = ?`

	return d.updateWholeId(query, a.FirstName, a.LastName, id)
}

func (d *Database) UpdateAuthor(id int64, a models.Author) error {
	fieldToDb := map[string]string{
		"FirstName": "imie",
		"LastName":  "nazwisko",
	}

	return d.updatePartId(a, "autor", id, fieldToDb)
}

func (d *Database) DelAuthor(id int64) error {
	query := "DELETE FROM autor WHERE id = ?"

	return d.deleteId(query, id)
}
