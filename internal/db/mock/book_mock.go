package mock

import (
	"net/url"

	"pawrest/internal/db"
	"pawrest/internal/models"
)

var books = []models.Book{
	{Id: 1, Title: "Book 1", Year: 1999, Pages: 300, Author: 1, Genre: 1, Language: 1},
	{Id: 2, Title: "Book 2", Year: 2005, Pages: 135, Author: 2, Genre: 1, Language: 2},
	{Id: 3, Title: "Book 3", Year: 1863, Pages: 48, Author: 3, Genre: 2, Language: 3},
}

var booksExt = []models.BookExt{
	{
		Id: 1, Title: "Book 1", Year: 1999, Pages: 300,
		Author:   models.Author{Id: 1, FirstName: "Adam", LastName: "Mickiewicz"},
		Genre:    models.Genre{Id: 1, Name: "Nowela"},
		Language: models.Language{Id: 1, Name: "Łaciński"},
	},
	{
		Id: 2, Title: "Book 2", Year: 2005, Pages: 135,
		Author:   models.Author{Id: 2, FirstName: "Witold", LastName: "Gombrowicz"},
		Genre:    models.Genre{Id: 1, Name: "Nowela"},
		Language: models.Language{Id: 2, Name: "Polski"},
	},
	{
		Id: 3, Title: "Book 3", Year: 1863, Pages: 48,
		Author:   models.Author{Id: 3, FirstName: "Bolesław", LastName: "Prus"},
		Genre:    models.Genre{Id: 2, Name: "Epopeja"},
		Language: models.Language{Id: 3, Name: "Angielski"},
	},
}

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

	return books, nil
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

	return booksExt, nil
}

func (m *MockDatabase) GetBook(id int64) (models.Book, error) {
	for _, book := range books {
		if book.Id == id {
			return book, nil
		}
	}

	return models.Book{}, db.ErrNotFound
}

func (m *MockDatabase) InsertBook(b models.Book) (int64, error) {
	if b.Language == 999 {
		return 0, db.ErrForeignKey
	}

	b.Id = int64(len(books) + 1)
	books = append(books, b)

	return b.Id, nil
}

func (m *MockDatabase) UpdateWholeBook(id int64, b models.Book) error {
	for i, book := range books {
		if b.Language == 999 {
			return db.ErrForeignKey
		}

		if book.Id == id {
			books[i] = b
			books[i].Id = id
			return nil
		}
	}

	return db.ErrNotFound
}

func (m *MockDatabase) UpdateBook(id int64, b models.Book) error {
	for i, book := range books {
		if b.Language == 999 {
			return db.ErrForeignKey
		}

		if book.Id == id {
			if b.Title != "" {
				books[i].Title = b.Title
			}
			if b.Year != 0 {
				books[i].Year = b.Year
			}
			if b.Pages > 0 {
				books[i].Pages = b.Pages
			}
			if b.Author > 0 {
				books[i].Author = b.Author
			}
			if b.Genre > 0 {
				books[i].Genre = b.Genre
			}
			if b.Language > 0 {
				books[i].Language = b.Language
			}

			return nil
		}
	}

	return db.ErrNotFound
}

func (m *MockDatabase) DelBook(id int64) error {
	for i, book := range books {
		if book.Id == id {
			books = append(books[:i], books[i+1:]...)
			return nil
		}
	}

	return db.ErrNotFound
}
