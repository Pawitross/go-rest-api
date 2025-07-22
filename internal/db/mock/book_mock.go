package mock

import (
	"net/url"

	"pawrest/internal/db"
	"pawrest/internal/models"
)

func (m *MockDatabase) GetBooks(params url.Values) ([]models.Book, error) {
	allowedParams := map[string]string{
		"id":       "id",
		"title":    "tytul",
		"year":     "rok_wydania",
		"pages":    "liczba_stron",
		"author":   "id_autora",
		"genre":    "id_gatunku",
		"language": "id_jezyka",
	}

	if len(params) > 0 {
		_, _, err := db.AssembleFilter(params, allowedParams)
		if err != nil {
			return []models.Book{}, err
		}
	}

	return m.Books, nil
}

func (m *MockDatabase) GetBooksExt(params url.Values) ([]models.BookExt, error) {
	allowedParams := map[string]string{
		"id":                "k.id",
		"title":             "tytul",
		"year":              "rok_wydania",
		"pages":             "liczba_stron",
		"author.id":         "id_autora",
		"author.first_name": "a.imie",
		"author.last_name":  "a.nazwisko",
		"author.birth_year": "a.rok_urodzenia",
		"author.death_year": "a.rok_smierci",
		"genre.id":          "id_gatunku",
		"genre.name":        "g.nazwa",
		"language.id":       "id_jezyka",
		"language.name":     "j.nazwa",
	}

	if len(params) > 0 {
		_, _, err := db.AssembleFilter(params, allowedParams)
		if err != nil {
			return []models.BookExt{}, err
		}
	}

	return m.BooksExt, nil
}

func (m *MockDatabase) GetBook(id int64) (models.Book, error) {
	for _, book := range m.Books {
		if book.ID == id {
			return book, nil
		}
	}

	return models.Book{}, db.ErrNotFound
}

func (m *MockDatabase) InsertBook(b models.Book) (int64, error) {
	if b.Language == 999 {
		return 0, db.ErrForeignKey
	}

	b.ID = int64(len(m.Books) + 1)
	m.Books = append(m.Books, b)

	return b.ID, nil
}

func (m *MockDatabase) UpdateWholeBook(id int64, b models.Book) error {
	for i, book := range m.Books {
		if b.Language == 999 {
			return db.ErrForeignKey
		}

		if book.ID == id {
			m.Books[i] = b
			m.Books[i].ID = id
			return nil
		}
	}

	return db.ErrNotFound
}

func (m *MockDatabase) UpdateBook(id int64, b models.Book) error {
	for i, book := range m.Books {
		if b.Language == 999 {
			return db.ErrForeignKey
		}

		if book.ID == id {
			if b.Title != "" {
				m.Books[i].Title = b.Title
			}
			if b.Year != 0 {
				m.Books[i].Year = b.Year
			}
			if b.Pages > 0 {
				m.Books[i].Pages = b.Pages
			}
			if b.Author > 0 {
				m.Books[i].Author = b.Author
			}
			if b.Genre > 0 {
				m.Books[i].Genre = b.Genre
			}
			if b.Language > 0 {
				m.Books[i].Language = b.Language
			}

			return nil
		}
	}

	return db.ErrNotFound
}

func (m *MockDatabase) DelBook(id int64) error {
	for i, book := range m.Books {
		if book.ID == id {
			m.Books = append(m.Books[:i], m.Books[i+1:]...)
			return nil
		}
	}

	return db.ErrNotFound
}
