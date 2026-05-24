package model

type Balance struct { // данные по балансу прльзователя
	UserID           string
	CurrentBalance   float64 `json:"current,omitempty"`
	WithdrawnBalance float64 `json:"withdrawn,omitempty"`
	UpdatedAt        int64
}

type TransactionWithdraw struct {
	OrderNumber string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}
