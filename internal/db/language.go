package db

import (
	"database/sql"
	"net/url"

	"pawrest/internal/models"
)

type LanguageDatabaseInterface interface {
	GetLanguages(params url.Values) ([]models.Language, error)
	GetLanguage(id int64) (models.Language, error)
	InsertLanguage(l models.Language) (int64, error)
	UpdateWholeLanguage(id int64, l models.Language) error
	DelLanguage(id int64) error
}

func (d *Database) GetLanguages(params url.Values) ([]models.Language, error) {
	query := `
	SELECT id, nazwa
	FROM jezyk`

	allowedParams := map[string]string{
		"id":   "id",
		"name": "nazwa",
	}

	langFunc := func(l *models.Language, rows *sql.Rows) error {
		return rows.Scan(&l.ID, &l.Name)
	}

	return queryWithParams[models.Language](
		d,
		query,
		params,
		allowedParams,
		langFunc,
	)
}

func (d *Database) GetLanguage(id int64) (models.Language, error) {
	query := `
	SELECT id, nazwa
	FROM jezyk
	WHERE id = ?`

	langFunc := func(l *models.Language, row *sql.Row) error {
		return row.Scan(&l.ID, &l.Name)
	}

	return queryID[models.Language](d, query, id, langFunc)
}

func (d *Database) InsertLanguage(l models.Language) (int64, error) {
	query := `
	INSERT INTO jezyk (nazwa)
	VALUES (?)`

	return d.insert(query, l.Name)
}

func (d *Database) UpdateWholeLanguage(id int64, l models.Language) error {
	query := `
	UPDATE jezyk
	SET
		nazwa = ?
	WHERE id = ?`

	return d.updateWholeID(query, l.Name, id)
}

func (d *Database) DelLanguage(id int64) error {
	query := "DELETE FROM jezyk WHERE id = ?"

	return d.deleteID(query, id)
}
