package models

type Book struct {
	Id      int64  `json:"id"`
	Tytul   string `json:"title"`
	Rok     int64  `json:"year"`
	Strony  int64  `json:"pages"`
	Autor   int64  `json:"author"`
	Gatunek int64  `json:"genre"`
	Jezyk   int64  `json:"language"`
}

func (b *Book) ValidateBook() bool {
	return b.Tytul == "" ||
		b.Rok == 0 ||
		b.Strony <= 0 ||
		b.Autor <= 0 ||
		b.Gatunek <= 0 ||
		b.Jezyk <= 0
}
