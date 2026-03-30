package model

type Balance struct { // данные по балансу прльзователя
	UserID            string
	CurrentBalance    float64 `json:"current,omitempty"`
	Withdrawn_balance float64 `json:"withdrawn,omitempty"`
	Updated_at        int64
}
