package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

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

	transaction, err := t.service.GetTransactionByID(transactionID)
	if err != nil {
		t.errorHandler("Error getting transaction: "+err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(transaction)
}

func (t *TransactionController) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	transactionDTO := t.decodeTransactionDTO(r)

	transaction, err := t.service.SaveTransaction(transactionDTO.ToTransaction())
	if err != nil {
		t.errorHandler("Error saving transaction: "+err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(transaction)
}

func (t *TransactionController) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	transactionID := t.validateTransactionID(r)

	transactionDTO := t.decodeTransactionDTO(r)

	transaction, err := t.service.UpdateTransactionByID(transactionID, transactionDTO.ToTransaction())
	if err != nil {
		t.errorHandler("Error updating transaction: "+err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(transaction)
}

func (t *TransactionController) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	transactionID := t.validateTransactionID(r)

	err := t.service.DeleteTransactionByID(transactionID)
	if err != nil {
		t.errorHandler("Error deleting transaction: "+err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (t *TransactionController) validateTransactionID(r *http.Request) int64 {
	params := mux.Vars(r)
	transactionIDPath := params["id"]
	if transactionIDPath == "" {
		t.errorHandler("Transaction ID is required", http.StatusBadRequest)
	}

	transactionID, err := strconv.ParseInt(transactionIDPath, 10, 64)
	if err != nil {
		t.errorHandler("Transaction ID must be a valid number: "+err.Error(), http.StatusBadRequest)
	}

	if transactionID <= 0 {
		t.errorHandler("Transaction ID must be a valid number", http.StatusBadRequest)
	}

	return transactionID
}

func (t *TransactionController) decodeTransactionDTO(r *http.Request) *presentation.TransactionDTO {
	var transactionDTO presentation.TransactionDTO

	if err := json.NewDecoder(r.Body).Decode(&transactionDTO); err != nil {
		t.errorHandler("Error decoding request body: "+err.Error(), http.StatusBadRequest)
	}

	if err := transactionDTO.ValidateRequest(); err != nil {
		t.errorHandler("Error validating request body: "+err.Error(), http.StatusBadRequest)
	}

	return &transactionDTO
}

func (t *TransactionController) errorHandler(errorMessage string, statusCode int) {
	panic(presentation.NewApiError(statusCode, errorMessage))
}
