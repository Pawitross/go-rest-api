package testutil

import (
	"database/sql"
	"log"
	"os"
	"strings"
)

func SetupDatabase(db *sql.DB) error {
	if err := setupFromScript(db); err != nil {
		log.Println("Error when setting up the database from the script:", err)

		log.Println("Setting up the database from hardcoded statements.")
		return setupHardcoded(db)
	}

	return nil
}

func setupFromScript(db *sql.DB) error {
	data, err := os.ReadFile("../../../sql/02-schema.sql")
	if err != nil {
		return err
	}

	splitData := strings.Split(string(data), ";")
	for _, statement := range splitData {
		st := strings.TrimSpace(statement)
		if st == "" {
			continue
		}

		if _, err := db.Exec(st); err != nil {
			return err
		}
	}

	return nil
}

func dropTables(db *sql.DB, tables []string) error {
	for _, table := range tables {
		_, err := db.Exec("DROP TABLE IF EXISTS " + table)
		if err != nil {
			return err
		}
	}

	return nil
}

func createTables(db *sql.DB, creates []string) error {
	for _, create := range creates {
		_, err := db.Exec("CREATE TABLE " + create)
		if err != nil {
			return err
		}
	}

	return nil
}

func populateTables(db *sql.DB, inserts []string) error {
	for _, insert := range inserts {
		_, err := db.Exec("INSERT INTO " + insert)
		if err != nil {
			return err
		}
	}

	return nil
}

func setupHardcoded(db *sql.DB) error {
	tables := []string{
		"ksiazka",
		"jezyk",
		"gatunek",
		"autor",
	}
	if err := dropTables(db, tables); err != nil {
		return err
	}

	creates := []string{
		`jezyk (
		    id      INT AUTO_INCREMENT,
		    nazwa   VARCHAR(64) NOT NULL,
		    PRIMARY KEY (id)
		);`,
		`gatunek (
		    id      INT AUTO_INCREMENT,
		    nazwa   VARCHAR(128) NOT NULL,
		    PRIMARY KEY (id)
		);`,
		`autor (
		    id              INT AUTO_INCREMENT,
		    imie            VARCHAR(128) NOT NULL,
		    nazwisko        VARCHAR(128) NOT NULL,
		    rok_urodzenia   DECIMAL(5) NOT NULL,
		    rok_smierci     DECIMAL(5),
		    PRIMARY KEY (id)
		);`,
		`ksiazka (
		    id              INT AUTO_INCREMENT,
		    tytul           VARCHAR(256) NOT NULL,
		    rok_wydania     DECIMAL(5) NOT NULL,
		    liczba_stron    INT,
		    okladka_url     VARCHAR(1024),
		    id_autora       INT NOT NULL,
		    id_gatunku      INT NOT NULL,
		    id_jezyka       INT NOT NULL,
		    PRIMARY KEY (id),
		    FOREIGN KEY (id_jezyka) REFERENCES jezyk(id),
		    FOREIGN KEY (id_autora) REFERENCES autor(id),
		    FOREIGN KEY (id_gatunku) REFERENCES gatunek(id)
		);`,
	}
	if err := createTables(db, creates); err != nil {
		return err
	}

	inserts := []string{
		`jezyk (nazwa) VALUES
		    ("Łaciński"),
		    ("Polski"),
		    ("Angielski"),
		    ("Niemiecki"),
		    ("Rosyjski"),
		    ("Francuski"),
		    ("Włoski"),
		    ("Hiszpański"),
		    ("Arabski"),
		    ("Chiński"),
		    ("Japoński");`,
		`gatunek (nazwa) VALUES
		    ("Nowela"),
		    ("Epopeja"),
		    ("Opowiadanie"),
		    ("Biografia"),
		    ("Dramat"),
		    ("Powieść"),
		    ("Opowieść"),
		    ("Zbiór poezji"),
		    ("Dystopia");`,
		`autor (imie, nazwisko, rok_urodzenia, rok_smierci) VALUES
		    ("Adam", "Mickiewicz", 1798, 1855),
		    ("Witold", "Gombrowicz", 1904, 1969),
		    ("Bolesław", "Prus", 1847, 1912),
		    ("Fiodor", "Dostojewski", 1821, 1881),
		    ("Stanisław", "Lem", 1921, 2006),
		    ("Jan", "Brzechwa", 1898, 1966),
		    ("Ernest", "Hemingway", 1899, 1961),
		    ("Henryk", "Sienkiewicz", 1846, 1916),
		    ("George", "Orwell", 1903, 1950);`,
		`ksiazka (
		    tytul, rok_wydania, liczba_stron, id_autora, id_gatunku, id_jezyka
		) VALUES
		    ("Pan Tadeusz, czyli ostatni zajazd na Litwie", 1834, 344, 1, 2, 2),
		    ("Dziady", 1822, 304, 1, 5, 2),
		    ("Ferdydurke", 1937, 296, 2, 6, 2),
		    ("Lalka", 1890, 676, 3, 6, 2),
		    ("Kamizelka", 1882, 24, 3, 1, 2),
		    ("Zbrodnia i kara", 1867, 496, 4, 6, 5),
		    ("Solaris", 1961, 340, 5, 6, 2),
		    ("Powrót z gwiazd", 1961, 400, 5, 6, 2),
		    ("Pokój na Ziemi", 1987, 376, 5, 6, 2),
		    ("Akademia pana Kleksa", 1946, 136, 6, 7, 2),
		    ("Brzechwa dzieciom", 1953, 176, 6, 8, 2),
		    ("Latarnik", 1881, 32, 8, 1, 2),
		    ("Ogniem i mieczem", 1884, 560, 8, 6, 2),
		    ("Potop", 1886, 936, 8, 6, 2),
		    ("Quo vadis", 1896, 448, 8, 6, 2),
		    ("Stary człowiek i morze", 1951, 100, 7, 3, 3),
		    ("Rok 1984", 1949, 312, 9, 9, 3);`,
	}
	if err := populateTables(db, inserts); err != nil {
		return err
	}

	return nil
}
