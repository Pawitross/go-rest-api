package sqldb

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var Db *sql.DB

type Ksiazka struct {
	Id      int64
	Tytul   string
	Rok     int64
	Autor   int64
	Gatunek int64
}

func ConnectToDB() {
	cfg := mysql.NewConfig()

	cfg.User = "root"
	cfg.Passwd = ""
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "paw"

	var err error
	Db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	if pingErr := Db.Ping(); pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Printf("Connected to database %q on %q", cfg.DBName, cfg.Addr)
}

func GetKsiazki() ([]Ksiazka, error) {
	var ksiazki []Ksiazka

	rows, err := Db.Query("SELECT * FROM Ksiazka")
	if err != nil {
		return nil, fmt.Errorf("GetKsiazki: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var ks Ksiazka

		if err := rows.Scan(&ks.Id, &ks.Tytul, &ks.Rok, &ks.Autor, &ks.Gatunek); err != nil {
			return nil, fmt.Errorf("GetKsiazki: %v", err)
		}

		ksiazki = append(ksiazki, ks)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetKsiazki: %v", err)
	}

	return ksiazki, nil
}

func GetKsiazka(id int64) (Ksiazka, error) {
	var ks Ksiazka

	row := Db.QueryRow("SELECT * FROM Ksiazka WHERE id = ?", id)
	if err := row.Scan(&ks.Id, &ks.Tytul, &ks.Rok, &ks.Autor, &ks.Gatunek); err != nil {
		if err == sql.ErrNoRows {
			return ks, fmt.Errorf("GetKsiazka: no book with id: %d", id)
		}

		return ks, fmt.Errorf("GetKsiazka: %v", err)
	}

	return ks, nil
}

func InsertKsiazka(ks Ksiazka) (int64, error) {
	result, err := Db.Exec("INSERT INTO ksiazka (tytul, rok_wydania, id_autora, id_gatunku) VALUES (?, ?, ?, ?)", ks.Tytul, ks.Rok, ks.Autor, ks.Gatunek)
	if err != nil {
		return 0, fmt.Errorf("InsertKsiazka: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("InsertKsiazka: %v", err)
	}

	return id, nil
}

func DelKsiazka(id int64) error {
	res, err := Db.Exec("DELETE FROM Ksiazka WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("DelKsiazka: %v", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("DelKsiazka, rows: %v", err)
	}

	if rows == 0 {
		return fmt.Errorf("DelKsiazka, no ksiazka with id: %v", id)
	}

	return nil
}
