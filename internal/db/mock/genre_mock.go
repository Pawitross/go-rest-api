package mock

import (
	"net/url"

	"pawrest/internal/db"
	"pawrest/internal/models"
)

func (m *MockDatabase) GetGenres(params url.Values) ([]models.Genre, error) {
	allowedParams := map[string]string{
		"id":   "id",
		"name": "nazwa",
	}

	if len(params) > 0 {
		_, _, err := db.AssembleFilter(params, allowedParams)
		if err != nil {
			return []models.Genre{}, err
		}
	}

	return m.Genres, nil
}

func (m *MockDatabase) GetGenre(id int64) (models.Genre, error) {
	for _, genre := range m.Genres {
		if genre.ID == id {
			return genre, nil
		}
	}

	return models.Genre{}, db.ErrNotFound
}

func (m *MockDatabase) InsertGenre(g models.Genre) (int64, error) {
	g.ID = int64(len(m.Genres) + 1)
	m.Genres = append(m.Genres, g)

	return g.ID, nil
}

func (m *MockDatabase) UpdateWholeGenre(id int64, g models.Genre) error {
	for i, genre := range m.Genres {
		if genre.ID == id {
			m.Genres[i] = g
			m.Genres[i].ID = id
			return nil
		}
	}

	return db.ErrNotFound
}

func (m *MockDatabase) DelGenre(id int64) error {
	for i, genre := range m.Genres {
		if genre.ID == id {
			m.Genres = append(m.Genres[:i], m.Genres[i+1:]...)
			return nil
		}
	}

	return db.ErrNotFound
}
