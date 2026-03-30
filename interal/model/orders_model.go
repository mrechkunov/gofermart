package model

type Orders struct { // все данные по заказу
	Number     int64  `json:"number" db:"o_number"`
	Status     string `json:"status" db:"o_status"`
	Accrual    int    `json:"accrual,omitempty" db:"o_accrual"`
	UploadedAt int64  `json:"uploaded_at" db:"uploaded_at"`
	CreatedBy  string `json:"created_by" db:"created_by"`
}

type AccrualOrder struct { // то что отдает сервис accrual
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

type ResponceOrders struct { // то что отдает наш сервис по запросу Get Orders
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}
