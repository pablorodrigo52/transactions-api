package service

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/pablorodrigo52/transaction-api/cmd/internal/presentation"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/repository"
)

type TransactionCurrencyService interface {
	GetTransactionCurrencyConverted(ctx context.Context, country string)
}

type TransactionCurrencyServiceImpl struct {
	treasuryRepository repository.TreasuryRepository
	log                *slog.Logger
}

func NewTransactionCurrencyService(treasuryRepository repository.TreasuryRepository, log *slog.Logger) *TransactionCurrencyServiceImpl {

	return &TransactionCurrencyServiceImpl{
		treasuryRepository: treasuryRepository,
		log:                log,
	}
}

func (s *TransactionCurrencyServiceImpl) GetTransactionCurrencyConverted(ctx context.Context, country string) {

	if country == "" {
		panic(presentation.NewApiError(http.StatusBadRequest, "invalid country name"))
	}

	// get treasury by country
	_, _ = s.treasuryRepository.GetExchangeRateByCountry(ctx, country)

	// validate if the effective date for the treasury is within in 6 months to the transactionDate
	// if yes then convert and return
	// if no then return an error stating the purchase cannot be converted to the target currency

	// need to return id, description, transactionDate, purchaseAmount, exchange rate, convertedPurchaseAmount
}
