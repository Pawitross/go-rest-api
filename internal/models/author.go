package models

type Author struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (a *Author) IsNotValid() bool {
	return a.FirstName == "" ||
		a.LastName == ""
}
