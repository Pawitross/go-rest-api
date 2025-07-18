package models

type Error struct {
	Error string `json:"error"`
} // @Name ErrorResponse

type Token struct {
	Admin bool   `json:"admin"`
	Token string `json:"token"`
} // @Name TokenResponse
