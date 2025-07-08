package models

type Genre struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (g *Genre) IsNotValid() bool {
	return g.Name == ""
}
