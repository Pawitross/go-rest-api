package mock

import (
	"net/url"

	"pawrest/internal/db"
	"pawrest/internal/models"
)

var genres = []models.Genre{
	{Id: 1, Name: "Science fiction"},
	{Id: 2, Name: "Dystopia"},
	{Id: 3, Name: "Biografia"},
	{Id: 4, Name: "Epopeja"},
	{Id: 5, Name: "Nowela"},
}

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

	return genres, nil
}

func (m *MockDatabase) GetGenre(id int64) (models.Genre, error) {
	for _, genre := range genres {
		if genre.Id == id {
			return genre, nil
		}
	}

	return models.Genre{}, db.ErrNotFound
}

func (m *MockDatabase) InsertGenre(a models.Genre) (int64, error) {
	a.Id = int64(len(genres) + 1)
	genres = append(genres, a)

	return a.Id, nil
}

func (m *MockDatabase) UpdateWholeGenre(id int64, a models.Genre) error {
	for i, genre := range genres {
		if genre.Id == id {
			genres[i] = a
			genres[i].Id = id
			return nil
		}
	}

	return db.ErrNotFound
}

func (m *MockDatabase) DelGenre(id int64) error {
	for i, genre := range genres {
		if genre.Id == id {
			genres = append(genres[:i], genres[i+1:]...)
			return nil
		}
	}

	return db.ErrNotFound
}
