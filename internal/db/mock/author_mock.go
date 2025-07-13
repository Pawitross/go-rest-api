package mock

import (
	"net/url"

	"pawrest/internal/db"
	"pawrest/internal/models"
)

func (m *MockDatabase) GetAuthors(params url.Values) ([]models.Author, error) {
	allowedParams := map[string]string{
		"id":         "id",
		"first_name": "imie",
		"last_name":  "nazwisko",
		"birth_year": "rok_urodzenia",
		"death_year": "rok_smierci",
	}

	if len(params) > 0 {
		_, _, err := db.AssembleFilter(params, allowedParams)
		if err != nil {
			return []models.Author{}, err
		}
	}

	return m.Authors, nil
}

func (m *MockDatabase) GetAuthor(id int64) (models.Author, error) {
	for _, author := range m.Authors {
		if author.Id == id {
			return author, nil
		}
	}

	return models.Author{}, db.ErrNotFound
}

func (m *MockDatabase) InsertAuthor(a models.Author) (int64, error) {
	a.Id = int64(len(m.Authors) + 1)
	m.Authors = append(m.Authors, a)

	return a.Id, nil
}

func (m *MockDatabase) UpdateWholeAuthor(id int64, a models.Author) error {
	for i, author := range m.Authors {
		if author.Id == id {
			m.Authors[i] = a
			m.Authors[i].Id = id
			return nil
		}
	}

	return db.ErrNotFound
}

func (m *MockDatabase) UpdateAuthor(id int64, a models.Author) error {
	for i, author := range m.Authors {
		if author.Id == id {
			if a.FirstName != "" {
				m.Authors[i].FirstName = a.FirstName
			}
			if a.LastName != "" {
				m.Authors[i].LastName = a.LastName
			}
			if a.BirthYear != 0 {
				m.Authors[i].BirthYear = a.BirthYear
			}
			if a.DeathYear != nil {
				m.Authors[i].DeathYear = a.DeathYear
			}

			return nil
		}
	}

	return db.ErrNotFound
}

func (m *MockDatabase) DelAuthor(id int64) error {
	for i, author := range m.Authors {
		if author.Id == id {
			m.Authors = append(m.Authors[:i], m.Authors[i+1:]...)
			return nil
		}
	}

	return db.ErrNotFound
}
