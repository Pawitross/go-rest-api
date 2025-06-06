package models

type Ksiazka struct {
	Id      int64  `json:"id"`
	Tytul   string `json:"title"`
	Rok     int64  `json:"year"`
	Strony  int64  `json:"pages"`
	Autor   int64  `json:"author"`
	Gatunek int64  `json:"genre"`
	Jezyk   int64  `json:"language"`
}
