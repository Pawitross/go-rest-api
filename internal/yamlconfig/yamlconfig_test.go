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
		t.Fatalf("Error creating file: %v", err)
	}
	defer os.Remove(fileName)
	f.Close()

	err = yamlconfig.Load(fileName)
	if err == nil {
		t.Errorf("Should return an error.")
	}

	expectedErr := "config file \"" + fileName + "\" is empty"
	if err.Error() != expectedErr {
		t.Errorf("Error loading file: %v", err)
	}
}

func TestLoad_YAMLUnmarshalError(t *testing.T) {
	fileName := "configtest.yaml"
	data := []byte(`Foo: "bar`)

	if err := os.WriteFile(fileName, data, 0644); err != nil {
		t.Fatalf("Error creating file: %v", err)
	}
	defer os.Remove(fileName)

	if err := yamlconfig.Load(fileName); err == nil {
		t.Errorf("Should return an error.")
	}
}

func TestLoad_Success_NoOutsideEnvSet(t *testing.T) {
	fileName := "configtest.yaml"
	data := []byte("Foo: \"bar\"\nBaz: \"qux\"")

	if err := os.WriteFile(fileName, data, 0644); err != nil {
		t.Fatalf("Error creating file: %v", err)
	}
	defer os.Remove(fileName)

	if err := yamlconfig.Load(fileName); err != nil {
		t.Fatalf("Error loading config: %v", err)
	}

	envVars := []struct {
		envVar  string
		wantVal string
	}{
		{"Foo", "bar"},
		{"Baz", "qux"},
	}

	for _, test := range envVars {
		val, found := os.LookupEnv(test.envVar)
		if !found {
			t.Errorf("Didn't find environment variable %v", test.envVar)
		}

		if val != test.wantVal {
			t.Errorf("Value of %v should be %v, not %v", test.envVar, test.wantVal, val)
		}
	}
}

func TestLoad_Success_OutsideEnvSet(t *testing.T) {
	os.Setenv("SETVAR", "foo")
	defer os.Unsetenv("SETVAR")

	fileName := "configtest.yaml"
	data := []byte("FILEVAR: \"bar\"\nSETVAR: \"file\"")

	if err := os.WriteFile(fileName, data, 0644); err != nil {
		t.Fatalf("Error creating file: %v", err)
	}
	defer os.Remove(fileName)

	if err := yamlconfig.Load(fileName); err != nil {
		t.Fatalf("Error loading config: %v", err)
	}

	envVars := []struct {
		envVar  string
		wantVal string
	}{
		{"FILEVAR", "bar"},
		{"SETVAR", "foo"},
	}

	for _, test := range envVars {
		val, found := os.LookupEnv(test.envVar)
		if !found {
			t.Errorf("Didn't find environment variable %v", test.envVar)
		}

		if val != test.wantVal {
			t.Errorf("Value of %v should be %v, not %v", test.envVar, test.wantVal, val)
		}
	}
}
