package mock

import (
	"net/url"

	"pawrest/internal/db"
	"pawrest/internal/models"
)

var languages = []models.Language{
	{Id: 1, Name: "Polski"},
	{Id: 2, Name: "Angielski"},
	{Id: 3, Name: "Łaciński"},
	{Id: 4, Name: "Niemiecki"},
	{Id: 5, Name: "Francuski"},
	{Id: 6, Name: "Rosyjski"},
}

func (m *MockDatabase) GetLanguages(params url.Values) ([]models.Language, error) {
	allowedParams := map[string]string{
		"id":   "id",
		"name": "nazwa",
	}

	if len(params) > 0 {
		_, _, err := db.AssembleFilter(params, allowedParams)
		if err != nil {
			return []models.Language{}, err
		}
	}

	return languages, nil
}

func (m *MockDatabase) GetLanguage(id int64) (models.Language, error) {
	for _, language := range languages {
		if language.Id == id {
			return language, nil
		}
	}

	return models.Language{}, db.ErrNotFound
}

func (m *MockDatabase) InsertLanguage(l models.Language) (int64, error) {
	l.Id = int64(len(languages) + 1)
	languages = append(languages, l)

	return l.Id, nil
}

func (m *MockDatabase) UpdateWholeLanguage(id int64, l models.Language) error {
	for i, language := range languages {
		if language.Id == id {
			languages[i] = l
			languages[i].Id = id
			return nil
		}
	}

	return db.ErrNotFound
}

func (m *MockDatabase) DelLanguage(id int64) error {
	for i, language := range languages {
		if language.Id == id {
			languages = append(languages[:i], languages[i+1:]...)
			return nil
		}
	}

	return db.ErrNotFound
}
