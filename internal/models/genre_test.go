package models_test

import (
	"testing"

	"pawrest/internal/models"
)

func TestGenreValidation(t *testing.T) {
	genreTests := map[string]struct {
		genre     models.Genre
		isInvalid bool
	}{
		"Valid": {
			genre: models.Genre{
				Id:   1,
				Name: "Horror",
			},
			isInvalid: false,
		},
		"InvalidEmptyName": {
			genre: models.Genre{
				Id:   1,
				Name: "",
			},
			isInvalid: true,
		},
	}

	for name, tt := range genreTests {
		t.Run(name, func(t *testing.T) {
			actual := tt.genre.IsNotValid()

			if tt.isInvalid != actual {
				t.Errorf("\n    test: %v\nexpected: %v\n     got: %v\n     for: %+v",
					name, tt.isInvalid, actual, tt.genre)
			}
		})
	}
}
