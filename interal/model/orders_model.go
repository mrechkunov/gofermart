package model

type Orders struct {
	Number     int64  `json:"number" db:"o_number"`
	Status     string `json:"status" db:"o_status"`
	Accrual    int    `json:"accrual,omitempty" db:"o_accrual"`
	UploadedAt int64  `json:"uploaded_at" db:"uploaded_at"`
	CreatedBy  string `json:"created_by" db:"created_by"`
}

type AccrualOrder struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual,omitempty"`
}
