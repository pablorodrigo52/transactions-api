package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/pablorodrigo52/transaction-api/cmd/internal/presentation"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/service"
)

type TransactionController struct {
	service service.TransactionService
	log     *slog.Logger
}

func NewTransactionController(log *slog.Logger, service service.TransactionService) *TransactionController {
	return &TransactionController{
		service: service,
		log:     log,
	}
}

func (t *TransactionController) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transactionDTO presentation.TransactionDTO

	if err := json.NewDecoder(r.Body).Decode(&transactionDTO); err != nil {
		t.errorHandler("Error decoding request body: "+err.Error(), http.StatusInternalServerError)
	}

	if err := transactionDTO.ValidateRequest(); err != nil {
		t.errorHandler("Error validating request body: "+err.Error(), http.StatusInternalServerError)
	}

	transaction, err := t.service.SaveTransaction(transactionDTO.ToTransaction())
	if err != nil {
		t.errorHandler("Error saving transaction: "+err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(transaction)
}

func (t *TransactionController) errorHandler(errorMessage string, statusCode int) {
	t.log.Error(errorMessage)
	panic(presentation.NewApiError(statusCode, errorMessage))
}
