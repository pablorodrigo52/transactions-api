package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/presentation"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/service"
)

type TransactionCurrencyController struct {
	service service.TransactionCurrencyService
	log     *slog.Logger
}

func NewTransactionCurrencyController(
	service service.TransactionCurrencyService,
	log *slog.Logger) *TransactionCurrencyController {

	return &TransactionCurrencyController{
		service: service,
		log:     log,
	}
}

func (c *TransactionCurrencyController) GetTransactionCurrency(w http.ResponseWriter, r *http.Request) {
	transactionID := c.validateTransactionID(r)
	country := c.validateCountryName(r)

	response := c.service.GetTransactionCurrencyConverted(r.Context(), transactionID, country)

	json.NewEncoder(w).Encode(response)
}

func (t *TransactionCurrencyController) validateTransactionID(r *http.Request) int64 {
	params := mux.Vars(r)
	transactionID := presentation.TransactionID(params["id"])

	transactionID.Validate()
	return transactionID.Get()
}

func (c *TransactionCurrencyController) validateCountryName(r *http.Request) string {
	params := mux.Vars(r)
	country := presentation.Country(params["country"])

	country.Validate()
	return country.Normalize()
}
