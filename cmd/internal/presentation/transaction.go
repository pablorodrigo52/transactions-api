package presentation

import (
	"net/http"

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

func (t *TransactionDTO) Validate() {
	if t.Description == "" || len(t.Description) > 50 {
		panic(NewApiError(http.StatusBadRequest, "invalid description, it must be between 1 and 50 characters"))
	}

	if t.TransactionDate == "" {
		panic(NewApiError(http.StatusBadRequest, "transaction date must not be empty"))
	}

	_, err := util.ParseDate(t.TransactionDate)
	if err != nil {
		panic(NewApiError(http.StatusBadRequest, err.Error()))
	}

	if t.PurchaseAmount <= 0 {
		panic(NewApiError(http.StatusBadRequest, "invalid purchase amount, it must be greater than 0"))
	}
}

func (t *TransactionDTO) ToTransaction() *model.Transaction {
	date, _ := util.ParseDate(t.TransactionDate)

	return &model.Transaction{
		ID:              t.TransactionID,
		Description:     t.Description,
		TransactionDate: date,
		PurchaseAmount:  util.RoundPurchaseAmount(t.PurchaseAmount),
	}
}
