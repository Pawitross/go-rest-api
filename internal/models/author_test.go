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
		"Valid": {
			author: models.Author{
				Id:        1,
				FirstName: "Jane",
				LastName:  "Doe",
			},
			isInvalid: false,
		},
		"EmptyFirstName": {
			author: models.Author{
				Id:        1,
				FirstName: "",
				LastName:  "Doe",
			},
			isInvalid: true,
		},
		"EmptyLastName": {
			author: models.Author{
				Id:        1,
				FirstName: "Jane",
				LastName:  "",
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
