package model

import (
	"time"
)

type Transaction struct {
	ID              int64
	Description     string
	TransactionDate time.Time
	PurchaseAmount  float32
	Deleted         bool
}
