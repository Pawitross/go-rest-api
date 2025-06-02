package sqldb

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-sql-driver/mysql"
)

var Db *sql.DB

type Ksiazka struct {
	Id      int64  `json:"id"`
	Tytul   string `json:"title"`
	Rok     int64  `json:"year"`
	Strony  int64  `json:"pages"`
	Autor   int64  `json:"author"`
	Gatunek int64  `json:"genre"`
	Jezyk   int64  `json:"language"`
}

func ConnectToDB() error {
	cfg := mysql.NewConfig()

	cfg.User = "root"
	cfg.Passwd = ""
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "paw"

	var err error
	Db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}

	if pingErr := Db.Ping(); pingErr != nil {
		return pingErr
	}

	fmt.Printf("Połączono z bazą danych %q on %q\n", cfg.DBName, cfg.Addr)
	return nil
}

func GetKsiazki(params url.Values) ([]Ksiazka, error) {
	query := "SELECT id, tytul, rok_wydania, liczba_stron, id_autora, id_gatunku, id_jezyka FROM Ksiazka"
	conditions := []string{}
	args := []any{}

	if len(params) != 0 {
		allowedParams := map[string]string{
			"id":       "id",
			"title":    "tytul",
			"year":     "rok_wydania",
			"pages":    "liczba_stron",
			"author":   "id_autora",
			"genre":    "id_gatunku",
			"language": "id_jezyka",
		}

		limit, limitted := params["limit"]
		if limitted {
			delete(params, "limit")
		}

		offset, offsetted := params["offset"]
		if offsetted {
			delete(params, "offset")
		}

		if !limitted && offsetted {
			return nil, fmt.Errorf("Nie podano limitu do podanego offsetu.")
		}

		for k, vsl := range params {
			dbCol, allowed := allowedParams[k]
			if !allowed || len(vsl) == 0 {
				return nil, fmt.Errorf("Wprowadzono nieznany parametr lub jest pusty.")
			}

			if len(vsl) > 1 {
				return nil, fmt.Errorf("Wprowadzono za dużo parametrów dla jednej kolumny.")
			}

			conditions = append(conditions, dbCol+" = ?")
			args = append(args, vsl[0])
		}

		if len(conditions) > 0 {
			query += " WHERE " + strings.Join(conditions, " AND ")
		}

		if limitted {
			query += " LIMIT ?"
			args = append(args, limit[0])
		}

		if offsetted {
			query += " OFFSET ?"
			args = append(args, offset[0])
		}
	}

	//fmt.Println("conditions:", conditions)
	//fmt.Println("query:", query)
	//fmt.Println("args:", args)
	//fmt.Println()

	var ksiazki []Ksiazka

	rows, err := Db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("Błąd zapytania (%v)", err)
	}
	defer rows.Close()

	for rows.Next() {
		var k Ksiazka

		if err := rows.Scan(&k.Id, &k.Tytul, &k.Rok, &k.Strony, &k.Autor, &k.Gatunek, &k.Jezyk); err != nil {
			return nil, fmt.Errorf("Błąd odczytywania (%v)", err)
		}

		ksiazki = append(ksiazki, k)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Błąd (%v)", err)
	}

	return ksiazki, nil
}

func ksiazkaExists(id int64) error {
	// TODO
	return fmt.Errorf("TODO")
}

func GetKsiazka(id int64) (Ksiazka, error) {
	var k Ksiazka

	query := "SELECT id, tytul, rok_wydania, liczba_stron, id_autora, id_gatunku, id_jezyka FROM Ksiazka WHERE id = ?"

	row := Db.QueryRow(query, id)
	if err := row.Scan(&k.Id, &k.Tytul, &k.Rok, &k.Strony, &k.Autor, &k.Gatunek, &k.Jezyk); err != nil {
		if err == sql.ErrNoRows {
			return k, fmt.Errorf("Brak książki o id %d", id)
		}

		return k, fmt.Errorf("Błąd odczytywania (%v)", err)
	}

	return k, nil
}

func InsertKsiazka(k Ksiazka) (int64, error) {
	query := "INSERT INTO ksiazka (tytul, rok_wydania, liczba_stron, id_autora, id_gatunku, id_jezyka) VALUES (?, ?, ?, ?, ?, ?)"

	result, err := Db.Exec(query, k.Tytul, k.Rok, k.Strony, k.Autor, k.Gatunek, k.Jezyk)
	if err != nil {
		return 0, fmt.Errorf("Nie udało się dodać rekordu (%v)", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Nie udało się pobrać id (%v)", err)
	}

	return id, nil
}

func UpdateWholeKsiazka(id int64, k Ksiazka) error {
	query := "UPDATE Ksiazka SET tytul = ?, rok_wydania = ?, liczba_stron = ?, id_autora = ?, id_gatunku = ?, id_jezyka = ? WHERE id = ?"

	res, err := Db.Exec(query, k.Tytul, k.Rok, k.Strony, k.Autor, k.Gatunek, k.Jezyk, id)
	if err != nil {
		return fmt.Errorf("Nie udało się zaktualizować (%v)", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Zmienione wiersze (%v)", err)
	}

	if rows == 0 {
		return fmt.Errorf("Nie znaleziono rekordu do aktualizacji lub nie zmieniono rekordu")
	}

	return nil
}

func DelKsiazka(id int64) error {
	query := "DELETE FROM Ksiazka WHERE id = ?"

	res, err := Db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Nie udało się usunąć (%v)", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Zmienione wiersze (%v)", err)
	}

	if rows == 0 {
		return fmt.Errorf("Brak ksiązki o id %v", id)
	}

	return nil
}
