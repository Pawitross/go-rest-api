package models

type Book struct {
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	Year     int64  `json:"year"`
	Pages    int64  `json:"pages"`
	Author   int64  `json:"author"`
	Genre    int64  `json:"genre"`
	Language int64  `json:"language"`
}

func (b *Book) IsNotValid() bool {
	return b.Title == "" ||
		b.Year == 0 ||
		b.Pages <= 0 ||
		b.Author <= 0 ||
		b.Genre <= 0 ||
		b.Language <= 0
}

type BookExt struct {
	Id       int64    `json:"id"`
	Title    string   `json:"title"`
	Year     int64    `json:"year"`
	Pages    int64    `json:"pages"`
	Author   Author   `json:"author"`
	Genre    Genre    `json:"genre"`
	Language Language `json:"language"`
}
