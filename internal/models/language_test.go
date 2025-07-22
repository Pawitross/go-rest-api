package models_test

import (
	"testing"

	"pawrest/internal/models"
)

func TestLanguageValidation(t *testing.T) {
	languageTests := map[string]struct {
		language  models.Language
		isInvalid bool
	}{
		"Valid": {
			language: models.Language{
				ID:   1,
				Name: "Horror",
			},
			isInvalid: false,
		},
		"InvalidEmptyName": {
			language: models.Language{
				ID:   1,
				Name: "",
			},
			isInvalid: true,
		},
	}

	for name, tt := range languageTests {
		t.Run(name, func(t *testing.T) {
			actual := tt.language.IsNotValid()

			if tt.isInvalid != actual {
				t.Errorf("\n    test: %v\nexpected: %v\n     got: %v\n     for: %+v",
					name, tt.isInvalid, actual, tt.language)
			}
		})
	}
}
