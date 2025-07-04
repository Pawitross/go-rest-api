package mock

import (
	"net/url"

	"pawrest/internal/db"
	"pawrest/internal/models"
)

var authors = []models.Author{
	{Id: 1, FirstName: "John", LastName: "Doe"},
	{Id: 2, FirstName: "Alice", LastName: "Smith"},
	{Id: 3, FirstName: "Richard", LastName: "Roe"},
}

func (m *MockDatabase) GetAuthors(params url.Values) ([]models.Author, error) {
	allowedParams := map[string]string{
		"id": "id",
	}

	if len(params) > 0 {
		_, _, err := db.AssembleFilter(params, allowedParams)
		if err != nil {
			return []models.Author{}, err
		}
	}

	return authors, nil
}

func (m *MockDatabase) GetAuthor(id int64) (models.Author, error) {
	for _, author := range authors {
		if author.Id == id {
			return author, nil
		}
	}

	return models.Author{}, db.ErrNotFound
}

func (m *MockDatabase) InsertAuthor(a models.Author) (int64, error) {
	a.Id = int64(len(authors) + 1)
	authors = append(authors, a)

	return a.Id, nil
}

func (m *MockDatabase) UpdateWholeAuthor(id int64, a models.Author) error {
	for i, author := range authors {
		if author.Id == id {
			authors[i] = a
			authors[i].Id = id
			return nil
		}
	}

	return db.ErrNotFound
}

func (m *MockDatabase) UpdateAuthor(id int64, a models.Author) error {
	for i, author := range authors {
		if author.Id == id {
			if a.FirstName != "" {
				authors[i].FirstName = a.FirstName
			}
			if a.LastName != "" {
				authors[i].LastName = a.LastName
			}

			return nil
		}
	}

	return db.ErrNotFound
}

func (m *MockDatabase) DelAuthor(id int64) error {
	for i, author := range authors {
		if author.Id == id {
			authors = append(authors[:i], authors[i+1:]...)
			return nil
		}
	}

	return db.ErrNotFound
}
