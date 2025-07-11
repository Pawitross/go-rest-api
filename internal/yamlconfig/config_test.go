package yamlconfig_test

import (
	"errors"
	"os"
	"testing"

	"pawrest/internal/yamlconfig"
)

func TestParse_Success(t *testing.T) {
	os.Clearenv()

	data := []byte(`DBNAME: "testdb"
DBUSER: "user"
DBHOST: "132.154.32.8"
DBPASS: "password"
SECRET: "secret-jwt-key"
`)

	fileName := "testenv.yaml"
	err := os.WriteFile(fileName, data, 0644)
	if err != nil {
		t.Fatalf("Error writing to file: %v", err)
	}
	defer os.Remove(fileName)

	cfg, err := yamlconfig.Parse(fileName)
	if err != nil {
		t.Fatalf("Should not return an error: %v", err)
	}

	values := []struct {
		want string
		got  string
	}{
		{"user", cfg.DBUser},
		{"password", cfg.DBPass},
		{"testdb", cfg.DBName},
		{"132.154.32.8", cfg.DBHost},
		{"3306", cfg.DBPort},
		{"secret-jwt-key", cfg.Secret},
	}

	for _, v := range values {
		if v.got != v.want {
			t.Errorf("got %v, want %v", v.got, v.want)
		}
	}
}

func TestParse_Error_MissingFile(t *testing.T) {
	_, err := yamlconfig.Parse("nonexist.yaml")
	if err == nil {
		t.Fatal("Shold return an error")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Error should be %v, not %v", "os.ErrNotExist", err.Error())
	}
}

func TestParse_Error_MissingVar(t *testing.T) {
	tests := map[string]struct {
		fileData      []byte
		wantErrSuffix string
	}{
		"DBUSER": {
			[]byte("DBNAME: \"testname\"\nSECRET: \"testsecret\""),
			"DBUSER",
		},
		"DBNAME": {
			[]byte("DBUSER: \"testuser\"\nSECRET: \"testsecret\""),
			"DBNAME",
		},
		"SECRET": {
			[]byte("DBUSER: \"testpass\"\nDBNAME: \"testname\""),
			"SECRET",
		},
		"DBUSER_DBNAME": {
			[]byte("SECRET: \"testsecret\""),
			"DBUSER, DBNAME",
		},
		"DBUSER_SECRET": {
			[]byte("DBNAME: \"testname\""),
			"DBUSER, SECRET",
		},
		"DBUSER_DBNAME_SECRET": {
			[]byte("ADDITIONAL: \"var\""),
			"DBUSER, DBNAME, SECRET",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			os.Clearenv()
			fileName := "testenv.yaml"

			err := os.WriteFile(fileName, tt.fileData, 0644)
			if err != nil {
				t.Fatalf("Error writing to file: %v", err)
			}
			defer os.Remove(fileName)

			_, err = yamlconfig.Parse(fileName)
			if err == nil {
				t.Fatal("Should return an error")
			}

			expectedErr := "Missing required environment variable/s: " + tt.wantErrSuffix
			if err.Error() != expectedErr {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
