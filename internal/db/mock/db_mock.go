package mock

import (
	"pawrest/internal/models"
)

type MockDatabase struct {
	Books     []models.Book
	BooksExt  []models.BookExt
	Authors   []models.Author
	Genres    []models.Genre
	Languages []models.Language
}

func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		Books: []models.Book{
			{ID: 1, Title: "Book 1", Year: 1999, Pages: 300, Author: 1, Genre: 1, Language: 1},
			{ID: 2, Title: "Book 2", Year: 2005, Pages: 135, Author: 2, Genre: 1, Language: 2},
			{ID: 3, Title: "Book 3", Year: 1863, Pages: 48, Author: 3, Genre: 2, Language: 3},
		},
		BooksExt: []models.BookExt{
			{
				ID: 1, Title: "Book 1", Year: 1999, Pages: 300,
				Author:   models.Author{ID: 1, FirstName: "Adam", LastName: "Mickiewicz"},
				Genre:    models.Genre{ID: 1, Name: "Nowela"},
				Language: models.Language{ID: 1, Name: "Łaciński"},
			},
			{
				ID: 2, Title: "Book 2", Year: 2005, Pages: 135,
				Author:   models.Author{ID: 2, FirstName: "Witold", LastName: "Gombrowicz"},
				Genre:    models.Genre{ID: 1, Name: "Nowela"},
				Language: models.Language{ID: 2, Name: "Polski"},
			},
			{
				ID: 3, Title: "Book 3", Year: 1863, Pages: 48,
				Author:   models.Author{ID: 3, FirstName: "Bolesław", LastName: "Prus"},
				Genre:    models.Genre{ID: 2, Name: "Epopeja"},
				Language: models.Language{ID: 3, Name: "Angielski"},
			},
		},
		Authors: []models.Author{
			{ID: 1, FirstName: "John", LastName: "Doe", BirthYear: 1949, DeathYear: models.I64Ptr(2023)},
			{ID: 2, FirstName: "Alice", LastName: "Smith", BirthYear: 1988, DeathYear: nil},
			{ID: 3, FirstName: "Richard", LastName: "Roe", BirthYear: 1921, DeathYear: models.I64Ptr(2009)},
		},
		Genres: []models.Genre{
			{ID: 1, Name: "Science fiction"},
			{ID: 2, Name: "Dystopia"},
			{ID: 3, Name: "Biografia"},
			{ID: 4, Name: "Epopeja"},
			{ID: 5, Name: "Nowela"},
		},
		Languages: []models.Language{
			{ID: 1, Name: "Polski"},
			{ID: 2, Name: "Angielski"},
			{ID: 3, Name: "Łaciński"},
			{ID: 4, Name: "Niemiecki"},
			{ID: 5, Name: "Francuski"},
			{ID: 6, Name: "Rosyjski"},
		},
	}
}
