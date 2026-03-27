package model

type Users struct {
	Login    string `json:"login" db:"ulogin"`
	Password string `json:"password" db:"upassword"`
	Bearer   string `db:"ubearer"`
}
