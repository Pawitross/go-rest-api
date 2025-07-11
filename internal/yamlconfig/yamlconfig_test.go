package yamlconfig_test

import (
	"errors"
	"os"
	"testing"

	"pawrest/internal/yamlconfig"
)

func TestLoad_NoConfigFile(t *testing.T) {
	fileName := "foo.yaml"

	err := yamlconfig.Load(fileName)
	if err == nil {
		t.Errorf("Should return an error.")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Error is not ErrNotExist: %v", err)
	}
}

func TestLoad_EmptyFile(t *testing.T) {
	fileName := "configtest.yaml"

	f, err := os.Create(fileName)
	if err != nil {
		t.Errorf("Error creating file: %v", err)
	}
	defer os.Remove(fileName)
	f.Close()

	err = yamlconfig.Load(fileName)
	if err == nil {
		t.Errorf("Should return an error.")
	}

	expectedErr := "Config file \"" + fileName + "\" is empty"
	if err.Error() != expectedErr {
		t.Errorf("Error loading file: %v", err)
	}
}

func TestLoad_YAMLUnmarshalError(t *testing.T) {
	fileName := "configtest.yaml"

	data := []byte(`Foo: "bar`)
	err := os.WriteFile(fileName, data, 0644)
	if err != nil {
		t.Errorf("Error creating file: %v", err)
	}
	defer os.Remove(fileName)

	if err := yamlconfig.Load(fileName); err == nil {
		t.Errorf("Should return an error.")
	}
}

func TestLoad_Success(t *testing.T) {
	fileName := "configtest.yaml"

	data := []byte("Foo: \"bar\"\nBaz: \"qux\"")
	err := os.WriteFile(fileName, data, 0644)
	if err != nil {
		t.Errorf("Error creating file: %v", err)
	}
	defer os.Remove(fileName)

	if err := yamlconfig.Load(fileName); err != nil {
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
