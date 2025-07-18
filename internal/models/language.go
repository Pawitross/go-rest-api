package models

type Language struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
} // @Name Language

func (l *Language) IsNotValid() bool {
	return l.Name == ""
}
