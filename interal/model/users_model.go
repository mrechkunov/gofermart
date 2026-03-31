package model

type Users struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Bearer   string
}
