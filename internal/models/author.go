package models

type Author struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthYear int64  `json:"birth_year"`
	DeathYear *int64 `json:"death_year"`
}

func (a *Author) IsNotValid() bool {
	return a.FirstName == "" ||
		a.LastName == "" ||
		(a.DeathYear != nil && (a.BirthYear > *a.DeathYear))
}

func I64Ptr(i int64) *int64 {
	return &i
}
