package model

import "time"

type Orders struct {
	Number     string    `json:"number" db:"o_number"`
	Status     string    `json:"status" db:"o_status"`
	Accrual    int       `json:"accrual,omitempty" db:"o_accrual"`
	UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
	CreatedBy  string    `json:"created_by" db:"created_by"`
}
