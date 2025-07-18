package models

type Genre struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
} // @Name Genre

func (g *Genre) IsNotValid() bool {
	return g.Name == ""
}
