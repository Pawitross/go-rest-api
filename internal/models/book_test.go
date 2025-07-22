package models_test

import (
	"testing"

	"pawrest/internal/models"
)

func modifyBook(book models.Book, modify func(b *models.Book)) models.Book {
	b := book
	modify(&b)
	return b
}

func TestBookValidation(t *testing.T) {
	validBook := models.Book{
		ID:       1,
		Title:    "Foo",
		Year:     1998,
		Pages:    100,
		Author:   1,
		Genre:    1,
		Language: 1,
	}

	bookTests := map[string]struct {
		book      models.Book
		isInvalid bool
	}{
		"Valid": {
			book:      validBook,
			isInvalid: false,
		},
		"InvalidEmptyTitle": {
			book: modifyBook(validBook, func(b *models.Book) {
				b.Title = ""
			}),
			isInvalid: true,
		},
		"InvalidZeroYear": {
			book: modifyBook(validBook, func(b *models.Book) {
				b.Year = 0
			}),
			isInvalid: true,
		},
		"InvalidNegativePages": {
			book: modifyBook(validBook, func(b *models.Book) {
				b.Pages = -1
			}),
			isInvalid: true,
		},
		"InvalidNegativeAuthor": {
			book: modifyBook(validBook, func(b *models.Book) {
				b.Author = -1
			}),
			isInvalid: true,
		},
		"InvalidNegativeGenre": {
			book: modifyBook(validBook, func(b *models.Book) {
				b.Genre = -1
			}),
			isInvalid: true,
		},
		"InvalidNegativeLanguage": {
			book: modifyBook(validBook, func(b *models.Book) {
				b.Language = -1
			}),
			isInvalid: true,
		},
	}

	for name, tt := range bookTests {
		t.Run(name, func(t *testing.T) {
			actual := tt.book.IsNotValid()

			if tt.isInvalid != actual {
				t.Errorf("\n    test: %v\nexpected: %v\n     got: %v\n     for: %+v",
					name, tt.isInvalid, actual, tt.book)
			}
		})
	}
}
