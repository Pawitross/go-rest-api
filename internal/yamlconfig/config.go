package yamlconfig

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DBUser string
	DBPass string
	DBName string
	DBHost string
	DBPort string
	Secret string
}

func Parse(fPath string) (*Config, error) {
	if err := Load(fPath); err != nil {
		return nil, fmt.Errorf("failed to load: %w", err)
	}

	var missing []string

	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbName := os.Getenv("DBNAME")
	if dbUser == "" {
		missing = append(missing, "DBUSER")
	}

	if dbName == "" {
		missing = append(missing, "DBNAME")
	}

	dbHost := os.Getenv("DBHOST")
	dbPort := os.Getenv("DBPORT")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}

	if dbPort == "" {
		dbPort = "3306"
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		missing = append(missing, "SECRET")
	}

	if len(missing) != 0 {
		return nil, fmt.Errorf("Missing required environment variable/s: %v", strings.Join(missing, ", "))
	}

	return &Config{
		DBUser: dbUser,
		DBPass: dbPass,
		DBName: dbName,
		DBHost: dbHost,
		DBPort: dbPort,
		Secret: secret,
	}, nil
}
