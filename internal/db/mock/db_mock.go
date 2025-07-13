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
			{Id: 1, Title: "Book 1", Year: 1999, Pages: 300, Author: 1, Genre: 1, Language: 1},
			{Id: 2, Title: "Book 2", Year: 2005, Pages: 135, Author: 2, Genre: 1, Language: 2},
			{Id: 3, Title: "Book 3", Year: 1863, Pages: 48, Author: 3, Genre: 2, Language: 3},
		},
		BooksExt: []models.BookExt{
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
		},
		Authors: []models.Author{
			{Id: 1, FirstName: "John", LastName: "Doe", BirthYear: 1949, DeathYear: models.I64Ptr(2023)},
			{Id: 2, FirstName: "Alice", LastName: "Smith", BirthYear: 1988, DeathYear: nil},
			{Id: 3, FirstName: "Richard", LastName: "Roe", BirthYear: 1921, DeathYear: models.I64Ptr(2009)},
		},
		Genres: []models.Genre{
			{Id: 1, Name: "Science fiction"},
			{Id: 2, Name: "Dystopia"},
			{Id: 3, Name: "Biografia"},
			{Id: 4, Name: "Epopeja"},
			{Id: 5, Name: "Nowela"},
		},
		Languages: []models.Language{
			{Id: 1, Name: "Polski"},
			{Id: 2, Name: "Angielski"},
			{Id: 3, Name: "Łaciński"},
			{Id: 4, Name: "Niemiecki"},
			{Id: 5, Name: "Francuski"},
			{Id: 6, Name: "Rosyjski"},
		},
	}
}
