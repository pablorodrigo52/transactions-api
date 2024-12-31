package presentation

import (
	"errors"
	"math"

	"github.com/pablorodrigo52/transaction-api/cmd/internal/model"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/util"
)

type TransactionDTO struct {
	TransactionID   int64   `json:"transaction_id"`
	Description     string  `json:"description"`
	TransactionDate string  `json:"transaction_date"`
	PurchaseAmount  float32 `json:"purchase_amount"`
	Deleted         bool    `json:"deleted,omitempty"`
}

func (t *TransactionDTO) ValidateRequest() error {
	if t.Description == "" || len(t.Description) > 50 {
		return errors.New("invalid description, it must be between 1 and 50 characters")
	}

	if t.TransactionDate == "" {
		return errors.New("transaction date must not be empty")
	}

	_, err := util.ParseDate(t.TransactionDate)
	if err != nil {
		return err
	}

	if t.PurchaseAmount <= 0 {
		return errors.New("invalid purchase amount, it must be greater than 0")
	}

	return nil
}

func (t *TransactionDTO) ToTransaction() *model.Transaction {
	date, _ := util.ParseDate(t.TransactionDate)

	return &model.Transaction{
		ID:              t.TransactionID,
		Description:     t.Description,
		TransactionDate: date,
		PurchaseAmount:  t.RoundPurchaseAmount(),
	}
}

func (t *TransactionDTO) RoundPurchaseAmount() float32 {
	return float32(math.Round(float64(t.PurchaseAmount*100)) / 100)
}
