package mock

import (
	"net/url"

	"pawrest/internal/db"
	"pawrest/internal/models"
)

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

	return m.Languages, nil
}

func (m *MockDatabase) GetLanguage(id int64) (models.Language, error) {
	for _, language := range m.Languages {
		if language.ID == id {
			return language, nil
		}
	}

	return models.Language{}, db.ErrNotFound
}

func (m *MockDatabase) InsertLanguage(l models.Language) (int64, error) {
	l.ID = int64(len(m.Languages) + 1)
	m.Languages = append(m.Languages, l)

	return l.ID, nil
}

func (m *MockDatabase) UpdateWholeLanguage(id int64, l models.Language) error {
	for i, language := range m.Languages {
		if language.ID == id {
			m.Languages[i] = l
			m.Languages[i].ID = id
			return nil
		}
	}

	return db.ErrNotFound
}

func (m *MockDatabase) DelLanguage(id int64) error {
	for i, language := range m.Languages {
		if language.ID == id {
			m.Languages = append(m.Languages[:i], m.Languages[i+1:]...)
			return nil
		}
	}

	return db.ErrNotFound
}
