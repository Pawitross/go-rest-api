package models

type Error struct {
	Error string `json:"error"`
}

type Token struct {
	Admin bool   `json:"admin"`
	Token string `json:"token"`
}
