package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
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

func (t *TransactionController) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	transactionID := t.validateTransactionID(r)

	transaction := t.service.GetTransactionByID(transactionID)

	json.NewEncoder(w).Encode(transaction)
}

func (t *TransactionController) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	transactionDTO := t.decodeTransactionDTO(r)

	transaction := t.service.SaveTransaction(transactionDTO.ToTransaction())

	json.NewEncoder(w).Encode(transaction)
}

func (t *TransactionController) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	transactionID := t.validateTransactionID(r)

	transactionDTO := t.decodeTransactionDTO(r)

	transaction := t.service.UpdateTransactionByID(transactionID, transactionDTO.ToTransaction())

	json.NewEncoder(w).Encode(transaction)
}

func (t *TransactionController) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	transactionID := t.validateTransactionID(r)

	t.service.DeleteTransactionByID(transactionID)

	w.WriteHeader(http.StatusNoContent)
}

func (t *TransactionController) validateTransactionID(r *http.Request) int64 {
	params := mux.Vars(r)
	transactionID := presentation.TransactionID(params["id"])

	transactionID.Validate()
	return transactionID.Get()
}

func (t *TransactionController) decodeTransactionDTO(r *http.Request) *presentation.TransactionDTO {
	var transactionDTO presentation.TransactionDTO

	if err := json.NewDecoder(r.Body).Decode(&transactionDTO); err != nil {
		t.errorHandler("Error decoding request body: "+err.Error(), http.StatusBadRequest)
	}

	transactionDTO.Validate()
	return &transactionDTO
}

func (t *TransactionController) errorHandler(errorMessage string, statusCode int) {
	panic(presentation.NewApiError(statusCode, errorMessage))
}
