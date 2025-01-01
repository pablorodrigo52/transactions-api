package service

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/model"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/presentation"
	mock_repository "github.com/pablorodrigo52/transaction-api/cmd/internal/repository/mocks"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/util"
	"github.com/stretchr/testify/assert"
)

func Test_GetTransactionCurrencyConverted(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	treasuryRepository := mock_repository.NewMockTreasuryRepository(mockCtrl)
	transactionRepository := mock_repository.NewMockTransactionRepository(mockCtrl)
	log := slog.Default()
	context := context.Background()

	service := NewTransactionCurrencyService(treasuryRepository, transactionRepository, log)

	t.Run("GetTransactionCurrencyConverted failed because invalid transaction id", func(t *testing.T) {
		// given
		transactionID := int64(0)
		country := ""
		expectedError := presentation.NewApiError(http.StatusBadRequest, "invalid transaction id")
		defer assertPanicErrors(t, expectedError)

		// when
		response := service.GetTransactionCurrencyConverted(context, transactionID, country)

		// then
		assert.Nil(t, response)
	})

	t.Run("GetTransactionCurrencyConverted failed because invalid country name", func(t *testing.T) {
		// given
		transactionID := int64(1)
		country := ""
		expectedError := presentation.NewApiError(http.StatusBadRequest, "invalid country name")
		defer assertPanicErrors(t, expectedError)

		// when
		response := service.GetTransactionCurrencyConverted(context, transactionID, country)

		// then
		assert.Nil(t, response)
	})

	t.Run("GetTransactionCurrencyConverted failed because transaction repository failed", func(t *testing.T) {
		// given
		transactionID := int64(1)
		country := "Brazil"
		errorMessage := "mock error"
		expectedError := presentation.NewApiError(http.StatusFailedDependency, errorMessage)
		defer assertPanicErrors(t, expectedError)

		transactionRepository.EXPECT().GetTransaction(transactionID).Return(nil, errors.New(errorMessage))

		// when
		response := service.GetTransactionCurrencyConverted(context, transactionID, country)

		// then
		assert.Nil(t, response)
	})

	t.Run("GetTransactionCurrencyConverted failed because transaction not found", func(t *testing.T) {
		// given
		transactionID := int64(1)
		country := "Brazil"
		errorMessage := "transaction not found"
		expectedError := presentation.NewApiError(http.StatusNotFound, errorMessage)
		defer assertPanicErrors(t, expectedError)

		transactionRepository.EXPECT().GetTransaction(transactionID).Return(nil, nil)

		// when
		response := service.GetTransactionCurrencyConverted(context, transactionID, country)

		// then
		assert.Nil(t, response)
	})

	t.Run("GetTransactionCurrencyConverted failed because treasury repository failed", func(t *testing.T) {
		// given
		transactionID := int64(1)
		country := "Brazil"
		errorMessage := "treasury repository error"
		expectedError := presentation.NewApiError(http.StatusBadGateway, errorMessage)
		defer assertPanicErrors(t, expectedError)

		transactionRepository.EXPECT().GetTransaction(transactionID).Return(&model.Transaction{}, nil)
		treasuryRepository.EXPECT().GetExchangeRateByCountry(context, country).Return(nil, errors.New(errorMessage))

		// when
		response := service.GetTransactionCurrencyConverted(context, transactionID, country)

		// then
		assert.Nil(t, response)
	})

	t.Run("GetTransactionCurrencyConverted failed because treasury repository returns no data", func(t *testing.T) {
		// given
		transactionID := int64(1)
		country := "Brazil"
		errorMessage := "purchase cannot be converted to the target currency: no data found"
		expectedError := presentation.NewApiError(http.StatusBadGateway, errorMessage)
		defer assertPanicErrors(t, expectedError)

		transactionRepository.EXPECT().GetTransaction(transactionID).Return(&model.Transaction{}, nil)
		treasuryRepository.EXPECT().GetExchangeRateByCountry(context, country).Return(&model.TreasuryRatesExchange{
			Data: []model.Data{},
		}, nil)

		// when
		response := service.GetTransactionCurrencyConverted(context, transactionID, country)

		// then
		assert.Nil(t, response)
	})

	t.Run("GetTransactionCurrencyConverted failed because treasury repository returns an invalid effective date format", func(t *testing.T) {
		// given
		transactionID := int64(1)
		country := "Brazil"
		errorMessage := "purchase cannot be converted to the target currency: not found effective rate to convert"
		expectedError := presentation.NewApiError(http.StatusBadGateway, errorMessage)
		defer assertPanicErrors(t, expectedError)

		transactionRepository.EXPECT().GetTransaction(transactionID).Return(&model.Transaction{
			TransactionDate: time.Now(),
		}, nil)
		treasuryRepository.EXPECT().GetExchangeRateByCountry(context, country).Return(&model.TreasuryRatesExchange{
			Data: []model.Data{
				{
					EffectiveDate: "2021-01-01T00:00:00Z",
				},
			}}, nil)

		// when
		response := service.GetTransactionCurrencyConverted(context, transactionID, country)

		// then
		assert.Nil(t, response)
	})

	t.Run("GetTransactionCurrencyConverted failed because not found effective rate to convert", func(t *testing.T) {
		// given
		transactionID := int64(1)
		country := "Brazil"
		errorMessage := "purchase cannot be converted to the target currency: not found effective rate to convert"
		expectedError := presentation.NewApiError(http.StatusBadGateway, errorMessage)
		defer assertPanicErrors(t, expectedError)

		transactionRepository.EXPECT().GetTransaction(transactionID).Return(&model.Transaction{
			TransactionDate: time.Now(),
		}, nil)
		treasuryRepository.EXPECT().GetExchangeRateByCountry(context, country).Return(&model.TreasuryRatesExchange{
			Data: []model.Data{
				{
					EffectiveDate: "2021-01-01",
				},
			}}, nil)

		// when
		response := service.GetTransactionCurrencyConverted(context, transactionID, country)

		// then
		assert.Nil(t, response)
	})

	t.Run("GetTransactionCurrencyConverted failed because invalid exchange rate", func(t *testing.T) {
		// given
		transactionID := int64(1)
		country := "Brazil"
		errorMessage := "purchase cannot be converted to the target currency: invalid exchange rate. rate=mock"
		expectedError := presentation.NewApiError(http.StatusBadGateway, errorMessage)
		defer assertPanicErrors(t, expectedError)

		transactionRepository.EXPECT().GetTransaction(transactionID).Return(&model.Transaction{
			TransactionDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		}, nil)
		treasuryRepository.EXPECT().GetExchangeRateByCountry(context, country).Return(&model.TreasuryRatesExchange{
			Data: []model.Data{
				{
					EffectiveDate: "2025-01-01",
					ExchangeRate:  "mock",
				},
			}}, nil)

		// when
		response := service.GetTransactionCurrencyConverted(context, transactionID, country)

		// then
		assert.Nil(t, response)
	})

	t.Run("GetTransactionCurrencyConverted with sucess", func(t *testing.T) {
		// given
		transactionID := int64(1)
		country := "Brazil"
		transaction := &model.Transaction{
			TransactionDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			PurchaseAmount:  1.74,
		}
		exchangeRate := &model.TreasuryRatesExchange{
			Data: []model.Data{
				{
					EffectiveDate: "2025-01-01",
					ExchangeRate:  "6.18",
				},
			}}
		expectedResponse := &presentation.TransactionCurrencyDTO{
			ConvertedPurchaseAmount: 10.75,
			PurchaseAmount:          transaction.PurchaseAmount,
			TransactionDate:         util.FormatDate(transaction.TransactionDate),
			ExchangeRate:            6.18,
		}

		transactionRepository.EXPECT().GetTransaction(transactionID).Return(transaction, nil)
		treasuryRepository.EXPECT().GetExchangeRateByCountry(context, country).Return(exchangeRate, nil)

		// when
		response := service.GetTransactionCurrencyConverted(context, transactionID, country)

		// then
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse, response)
	})
}

func assertPanicErrors(t *testing.T, expectedError *presentation.ApiError) {
	if r := recover(); r != nil {
		assert.Equal(t, expectedError, r)
	}
}
