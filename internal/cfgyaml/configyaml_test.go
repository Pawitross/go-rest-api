package cfgyaml_test

import (
	"errors"
	"os"
	"testing"

	"pawrest/internal/cfgyaml"
)

func TestNoConfigFile(t *testing.T) {
	fileName := "foo.yaml"

	err := cfgyaml.Load(fileName)
	if err == nil {
		t.Errorf("Should return an error.")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Error is not ErrNotExist: %v", err)
	}
}

func TestEmptyFile(t *testing.T) {
	fileName := "configtest.yaml"
	defer os.Remove(fileName)

	f, err := os.Create(fileName)
	if err != nil {
		t.Errorf("Error creating file: %v", err)
	}
	f.Close()

	err = cfgyaml.Load(fileName)
	if err == nil {
		t.Errorf("Should return an error.")
	}

	expectedErr := "Config file \"" + fileName + "\" is empty."
	if err.Error() != expectedErr {
		t.Errorf("Error loading file: %v", err)
	}
}

func TestYAMLUnmarshalError(t *testing.T) {
	fileName := "configtest.yaml"
	defer os.Remove(fileName)

	f, err := os.Create(fileName)
	if err != nil {
		t.Errorf("Error creating file: %v", err)
	}

	f.Write([]byte(`Foo: "bar`))
	f.Close()

	if err := cfgyaml.Load(fileName); err == nil {
		t.Errorf("Should return an error.")
	}
}

func TestSuccess(t *testing.T) {
	fileName := "configtest.yaml"
	defer os.Remove(fileName)

	f, err := os.Create(fileName)
	if err != nil {
		t.Errorf("Error creating file: %v", err)
	}

	f.Write([]byte(`Foo: "bar"`))
	f.Write([]byte("\n"))
	f.Write([]byte(`Baz: "qux"`))
	f.Close()

	if err := cfgyaml.Load(fileName); err != nil {
		t.Fatalf("Error loading config: %v", err)
	}

	valFoo, foundFoo := os.LookupEnv("Foo")
	if !foundFoo {
		t.Errorf("Didn't found environment variable")
	}

	if valFoo != "bar" {
		t.Errorf(`Value of Foo should be "bar"`)
	}

	valBaz, foundBaz := os.LookupEnv("Baz")
	if !foundBaz {
		t.Errorf("Didn't found environment variable")
	}

	if valBaz != "qux" {
		t.Errorf(`Value of baz should be "qux"`)
	}
}
