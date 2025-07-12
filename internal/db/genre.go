package db

import (
	"database/sql"
	"net/url"

	"pawrest/internal/models"
)

type GenreDatabaseInterface interface {
	GetGenres(params url.Values) ([]models.Genre, error)
	GetGenre(id int64) (models.Genre, error)
	InsertGenre(g models.Genre) (int64, error)
	UpdateWholeGenre(id int64, g models.Genre) error
	DelGenre(id int64) error
}

func (d *Database) GetGenres(params url.Values) ([]models.Genre, error) {
	query := `
	SELECT id, nazwa
	FROM gatunek`

	allowedParams := map[string]string{
		"id":   "id",
		"name": "nazwa",
	}

	genreFunc := func(g *models.Genre, rows *sql.Rows) error {
		return rows.Scan(&g.Id, &g.Name)
	}

	return queryWithParams[models.Genre](
		d,
		query,
		params,
		allowedParams,
		genreFunc,
	)
}

func (d *Database) GetGenre(id int64) (models.Genre, error) {
	query := `
	SELECT id, nazwa
	FROM gatunek
	WHERE id = ?`

	genreFunc := func(g *models.Genre, row *sql.Row) error {
		return row.Scan(&g.Id, &g.Name)
	}

	return queryId[models.Genre](d, query, id, genreFunc)
}

func (d *Database) InsertGenre(g models.Genre) (int64, error) {
	query := `
	INSERT INTO gatunek (nazwa)
	VALUES (?)`

	return d.insert(query, g.Name)
}

func (d *Database) UpdateWholeGenre(id int64, g models.Genre) error {
	query := `
	UPDATE gatunek
	SET
		nazwa = ?
	WHERE id = ?`

	return d.updateWholeId(query, g.Name, id)
}

func (d *Database) DelGenre(id int64) error {
	query := "DELETE FROM gatunek WHERE id = ?"

	return d.deleteId(query, id)
}
