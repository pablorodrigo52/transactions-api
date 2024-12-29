package presentation

import (
	"errors"

	"github.com/pablorodrigo52/transaction-api/cmd/internal/model"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/util"
)

type TransactionDTO struct {
	TransactionID   int64   `json:"transaction_id"`
	Description     string  `json:"description"`
	TransactionDate string  `json:"transaction_date"`
	PurchaseAmount  float32 `json:"purchase_amount"`
}

func (t *TransactionDTO) ValidateRequest() error {
	if t.Description == "" || len(t.Description) > 50 {
		return errors.New("invalid description, it must be between 1 and 50 characters")
	}

	_, err := util.ParseDate(t.TransactionDate)
	if t.TransactionDate == "" || err != nil {
		return errors.New("invalid transaction date, the correct format is YYYY-MM-DD")
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
		PurchaseAmount:  t.PurchaseAmount,
	}
}
