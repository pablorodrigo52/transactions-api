package service

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/pablorodrigo52/transaction-api/cmd/internal/model"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/presentation"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/repository"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/util"
)

const (
	exchangeRateDateFormat = "2006-01-02"
)

type TransactionCurrencyService interface {
	GetTransactionCurrencyConverted(ctx context.Context, transactionID int64, country string) *presentation.TransactionCurrencyDTO
}

//go:generate mockgen -source=./transaction_currency_service.go -destination=./mocks/transaction_currency_service_mock.go

type TransactionCurrencyServiceImpl struct {
	treasuryRepository    repository.TreasuryRepository
	transactionRepository repository.TransactionRepository
	log                   *slog.Logger
}

func NewTransactionCurrencyService(treasuryRepository repository.TreasuryRepository, transactionRepository repository.TransactionRepository, log *slog.Logger) *TransactionCurrencyServiceImpl {

	return &TransactionCurrencyServiceImpl{
		treasuryRepository:    treasuryRepository,
		transactionRepository: transactionRepository,
		log:                   log,
	}
}

func (s *TransactionCurrencyServiceImpl) GetTransactionCurrencyConverted(ctx context.Context, transactionID int64, country string) *presentation.TransactionCurrencyDTO {

	if transactionID <= 0 {
		s.throwError(http.StatusBadRequest, "invalid transaction id")
	}

	if country == "" {
		panic(presentation.NewApiError(http.StatusBadRequest, "invalid country name"))
	}

	// get transaction by id
	trx, err := s.transactionRepository.GetTransaction(transactionID)
	if err != nil {
		s.throwError(http.StatusFailedDependency, err.Error())
	}

	if trx == nil || trx.Deleted {
		s.throwError(http.StatusNotFound, "transaction not found")
	}

	// get treasury by country
	exchangeRate, err := s.treasuryRepository.GetExchangeRateByCountry(ctx, country)
	if err != nil {
		s.throwError(http.StatusBadGateway, err.Error())
	}

	if len(exchangeRate.Data) == 0 {
		s.throwError(http.StatusBadGateway, "purchase cannot be converted to the target currency: no data found")
	}

	if !s.isAbleToConvertToTargetCurrency(trx.TransactionDate, *exchangeRate) {
		s.throwError(http.StatusBadGateway, "purchase cannot be converted to the target currency: not found effective rate to convert")
	}

	exchangeRateConverted, err := strconv.ParseFloat(exchangeRate.Data[0].ExchangeRate, 32)
	if err != nil {
		s.throwError(http.StatusBadGateway, "purchase cannot be converted to the target currency: invalid exchange rate. rate="+exchangeRate.Data[0].ExchangeRate)
	}

	return &presentation.TransactionCurrencyDTO{
		TransactionID:           trx.ID,
		Description:             trx.Description,
		TransactionDate:         util.FormatDate(trx.TransactionDate),
		PurchaseAmount:          trx.PurchaseAmount,
		ExchangeRate:            float32(exchangeRateConverted),
		ConvertedPurchaseAmount: util.RoundPurchaseAmount(trx.PurchaseAmount * float32(exchangeRateConverted)),
	}
}

// isAbleToConvertToTargetCurrency validates if the transaction date is within 6 months of the effective rate date
func (s *TransactionCurrencyServiceImpl) isAbleToConvertToTargetCurrency(transactionDate time.Time, exchangeRate model.TreasuryRatesExchange) bool {
	effectiveDateParsed, err := util.ParseDateWithFormat(exchangeRate.Data[0].EffectiveDate, exchangeRateDateFormat)
	if err != nil {
		return false
	}

	return effectiveDateParsed.
		AddDate(0, 6, 0).
		Compare(transactionDate) >= 0
}

func (s *TransactionCurrencyServiceImpl) throwError(status int, message string) {
	panic(presentation.NewApiError(status, message))
}
