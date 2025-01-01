package controller

import (
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
	// get transaction id + validate transaction id
	// transactionID := t.validateTransactionID(r)

	// get country name + validate country name

	// call service

	// return json encode response
}

// func (t *TransactionCurrencyController) validateTransactionID(r *http.Request) int64 {
// 	params := mux.Vars(r)
// 	transactionIDPath := params["id"]
// 	if transactionIDPath == "" {
// 		t.errorHandler("Transaction ID is required", http.StatusBadRequest)
// 	}

// 	transactionID, err := strconv.ParseInt(transactionIDPath, 10, 64)
// 	if err != nil {
// 		t.errorHandler("Transaction ID must be a valid number: "+err.Error(), http.StatusBadRequest)
// 	}

// 	if transactionID <= 0 {
// 		t.errorHandler("Transaction ID must be a valid number", http.StatusBadRequest)
// 	}

// 	return transactionID
// }

// func (t *TransactionCurrencyController) validateCountryName(r *http.Request) string {
// 	params := mux.Vars(r)
// 	countryName := params["country"]

// 	var country presentation.Country
// 	country = countryName

// 	return country.Normalize()
// }

func (c *TransactionCurrencyController) validateCountryName(r *http.Request) string {
	params := mux.Vars(r)
	country := presentation.Country(params["country"])

	country.Validate()
	return country.Normalize()
}
