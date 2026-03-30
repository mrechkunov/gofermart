package model

type Balance struct { // данные по балансу прльзователя
	UserID            string
	CurrentBalance    float64 `json:"current,omitempty"`
	Withdrawn_balance float64 `json:"withdrawn,omitempty"`
	Updated_at        int64
}

type TransactionWithdraw struct {
	OrderNumber  string  `json:"order"`
	Sum          float64 `json:"sum"`
	Processed_at string  `json:"processed_at"`
}
