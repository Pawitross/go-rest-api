package models_test

import (
	"testing"

	"pawrest/internal/models"
)

func TestAuthorValidation(t *testing.T) {
	authorTests := map[string]struct {
		author    models.Author
		isInvalid bool
	}{
		"ValidFull": {
			author: models.Author{
				ID:        1,
				FirstName: "Jane",
				LastName:  "Doe",
				BirthYear: 1970,
				DeathYear: models.I64Ptr(2043),
			},
			isInvalid: false,
		},
		"ValidNilDeathYear": {
			author: models.Author{
				ID:        1,
				FirstName: "Jane",
				LastName:  "Doe",
				BirthYear: 1970,
				DeathYear: nil,
			},
			isInvalid: false,
		},
		"InvalidEmptyFirstName": {
			author: models.Author{
				ID:        1,
				FirstName: "",
				LastName:  "Doe",
				BirthYear: 1970,
				DeathYear: models.I64Ptr(2043),
			},
			isInvalid: true,
		},
		"InvalidEmptyLastName": {
			author: models.Author{
				ID:        1,
				FirstName: "Jane",
				LastName:  "",
				BirthYear: 1970,
				DeathYear: models.I64Ptr(2043),
			},
			isInvalid: true,
		},
		"InvalidBirthYearGreaterThanDeathYear": {
			author: models.Author{
				ID:        1,
				FirstName: "Jane",
				LastName:  "Doe",
				BirthYear: 2000,
				DeathYear: models.I64Ptr(1900),
			},
			isInvalid: true,
		},
	}

	for name, tt := range authorTests {
		t.Run(name, func(t *testing.T) {
			actual := tt.author.IsNotValid()

			if tt.isInvalid != actual {
				t.Errorf("\n    test: %v\nexpected: %v\n     got: %v\n     for: %+v",
					name, tt.isInvalid, actual, tt.author)
			}
		})
	}
}
