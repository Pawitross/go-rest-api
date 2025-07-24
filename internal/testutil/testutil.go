package testutil

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

func SetupDatabase(db *sql.DB) error {
	if err := setupFromScript(db); err != nil {
		return fmt.Errorf("failed to set up the database from script: %w", err)
	}

	return nil
}

func setupFromScript(db *sql.DB) error {
	data, err := os.ReadFile("../../../sql/02-schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	splitData := strings.Split(string(data), ";")
	for _, statement := range splitData {
		st := strings.TrimSpace(statement)
		if st == "" {
			continue
		}

		if _, err := db.Exec(st); err != nil {
			return fmt.Errorf("failed to exec: %w", err)
		}
	}

	return nil
}
