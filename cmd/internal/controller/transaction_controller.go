package controller

import (
	"encoding/json"
	"errors"
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
	transactionID, err := t.validateTransactionID(r)
	if err != nil {
		return
	}

	transaction, err := t.service.GetTransactionByID(transactionID)
	if err != nil {
		t.errorHandler("Error getting transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transaction)
}

func (t *TransactionController) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	transactionDTO, err := t.decodeTransactionDTO(r)
	if err != nil {
		return
	}

	transaction, err := t.service.SaveTransaction(transactionDTO.ToTransaction())
	if err != nil {
		t.errorHandler("Error saving transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transaction)
}

func (t *TransactionController) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	transactionID, err := t.validateTransactionID(r)
	if err != nil {
		return
	}

	transactionDTO, err := t.decodeTransactionDTO(r)
	if err != nil {
		return
	}

	transaction, err := t.service.UpdateTransactionByID(transactionID, transactionDTO.ToTransaction())
	if err != nil {
		t.errorHandler("Error updating transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transaction)
}

func (t *TransactionController) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	transactionID, err := t.validateTransactionID(r)
	if err != nil {
		return
	}

	err = t.service.DeleteTransactionByID(transactionID)
	if err != nil {
		t.errorHandler("Error deleting transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (t *TransactionController) validateTransactionID(r *http.Request) (int64, error) {
	params := mux.Vars(r)
	transactionIDPath := params["id"]
	if transactionIDPath == "" {
		t.errorHandler("Transaction ID is required", http.StatusBadRequest)
		return 0, errors.New("transaction ID is required")
	}

	transactionID, err := strconv.ParseInt(transactionIDPath, 10, 64)
	if err != nil {
		t.errorHandler("Transaction ID must be a valid number: "+err.Error(), http.StatusBadRequest)
		return 0, err
	}

	if transactionID <= 0 {
		t.errorHandler("Transaction ID must be a valid number", http.StatusBadRequest)
		return 0, errors.New("transaction ID must be a valid number")
	}

	return transactionID, nil
}

func (t *TransactionController) decodeTransactionDTO(r *http.Request) (*presentation.TransactionDTO, error) {
	var transactionDTO presentation.TransactionDTO

	if err := json.NewDecoder(r.Body).Decode(&transactionDTO); err != nil {
		t.errorHandler("Error decoding request body: "+err.Error(), http.StatusBadRequest)
		return nil, err
	}

	if err := transactionDTO.ValidateRequest(); err != nil {
		t.errorHandler("Error validating request body: "+err.Error(), http.StatusBadRequest)
		return nil, err
	}

	return &transactionDTO, nil
}

func (t *TransactionController) errorHandler(errorMessage string, statusCode int) {
	panic(presentation.NewApiError(statusCode, errorMessage))
}
