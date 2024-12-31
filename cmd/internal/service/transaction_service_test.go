package service

import (
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

func Test_TransactionService_GetTransactionByID(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	mockRepository := mock_repository.NewMockTransactionRepository(mockController)
	mockCache := mock_repository.NewMockTransactionCache(mockController)

	transactionService := NewTransactionService(slog.Default(), mockRepository, mockCache)

	t.Run("Get transaction by id with success", func(t *testing.T) {
		// given
		mockTransaction := model.Transaction{
			ID:              int64(1),
			Description:     "mock description",
			TransactionDate: time.Now(),
			PurchaseAmount:  1.0,
		}

		// when
		mockCache.EXPECT().Get(mockTransaction.ID).Return(nil)
		mockCache.EXPECT().Save(mockTransaction.ID, &mockTransaction).Return(nil)
		mockRepository.EXPECT().GetTransaction(mockTransaction.ID).Return(&mockTransaction, nil)

		response, err := transactionService.GetTransactionByID(mockTransaction.ID)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, response)

		expectedTransaction := &presentation.TransactionDTO{
			TransactionID:   mockTransaction.ID,
			Description:     mockTransaction.Description,
			TransactionDate: util.FormatDate(mockTransaction.TransactionDate),
			PurchaseAmount:  mockTransaction.PurchaseAmount,
		}
		assert.Equal(t, expectedTransaction, response)
	})
	t.Run("Get transaction by id with success recover from cache", func(t *testing.T) {
		// given
		mockTransaction := model.Transaction{
			ID:              int64(1),
			Description:     "mock description",
			TransactionDate: time.Now(),
			PurchaseAmount:  1.0,
		}

		// when
		mockCache.EXPECT().Get(mockTransaction.ID).Return(&mockTransaction)
		response, err := transactionService.GetTransactionByID(mockTransaction.ID)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, response)

		expectedTransaction := &presentation.TransactionDTO{
			TransactionID:   mockTransaction.ID,
			Description:     mockTransaction.Description,
			TransactionDate: util.FormatDate(mockTransaction.TransactionDate),
			PurchaseAmount:  mockTransaction.PurchaseAmount,
		}
		assert.Equal(t, expectedTransaction, response)
	})
	t.Run("Get transaction by id with success but error on save cache", func(t *testing.T) {
		// given
		mockTransaction := model.Transaction{
			ID:              int64(1),
			Description:     "mock description",
			TransactionDate: time.Now(),
			PurchaseAmount:  1.0,
		}

		// when
		mockCache.EXPECT().Get(mockTransaction.ID).Return(nil)
		mockCache.EXPECT().Save(mockTransaction.ID, &mockTransaction).Return(errors.New("mock error"))
		mockRepository.EXPECT().GetTransaction(mockTransaction.ID).Return(&mockTransaction, nil)

		response, err := transactionService.GetTransactionByID(mockTransaction.ID)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, response)

		expectedTransaction := &presentation.TransactionDTO{
			TransactionID:   mockTransaction.ID,
			Description:     mockTransaction.Description,
			TransactionDate: util.FormatDate(mockTransaction.TransactionDate),
			PurchaseAmount:  mockTransaction.PurchaseAmount,
		}
		assert.Equal(t, expectedTransaction, response)
	})
	t.Run("Get transaction by id error transaction not found", func(t *testing.T) {
		// given
		mockTransaction := model.Transaction{
			ID: int64(1),
		}
		expectedError := presentation.NewApiError(http.StatusNotFound, "transaction not found")

		// when
		mockCache.EXPECT().Get(mockTransaction.ID).Return(nil)
		mockRepository.EXPECT().GetTransaction(mockTransaction.ID).Return(nil, nil)

		// then
		defer assertPanicApiErrors(t, expectedError)

		_, _ = transactionService.GetTransactionByID(mockTransaction.ID)
	})
	t.Run("Get transaction by id error getting transaction", func(t *testing.T) {
		// given
		mockTransaction := model.Transaction{
			ID: int64(1),
		}
		expectedError := presentation.NewApiError(http.StatusInternalServerError, "error getting transaction")

		// when
		mockCache.EXPECT().Get(mockTransaction.ID).Return(nil)
		mockRepository.EXPECT().GetTransaction(mockTransaction.ID).Return(nil, errors.New("mock error"))

		// then
		defer assertPanicApiErrors(t, expectedError)

		_, _ = transactionService.GetTransactionByID(mockTransaction.ID)
	})
	t.Run("Get transaction by id error invalid transaction id", func(t *testing.T) {
		// given
		mockTransaction := model.Transaction{
			ID: int64(0),
		}
		expectedError := presentation.NewApiError(http.StatusBadRequest, "invalid transaction id: 0")

		// then
		defer assertPanicApiErrors(t, expectedError)

		// when
		_, _ = transactionService.GetTransactionByID(mockTransaction.ID)
	})
}

func Test_TransactionService_SaveTransaction(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	mockRepository := mock_repository.NewMockTransactionRepository(mockController)
	mockCache := mock_repository.NewMockTransactionCache(mockController)

	transactionService := NewTransactionService(slog.Default(), mockRepository, mockCache)

	t.Run("Save transaction with success", func(t *testing.T) {
		// given
		mockTransaction := model.Transaction{
			Description:     "mock description",
			TransactionDate: time.Now(),
			PurchaseAmount:  1.0,
		}
		savedTransaction := mockTransaction
		savedTransaction.ID = int64(1)

		// when
		mockRepository.EXPECT().SaveTransaction(&mockTransaction).Return(&savedTransaction, nil)
		mockCache.EXPECT().Save(savedTransaction.ID, &savedTransaction).Return(nil)

		response, err := transactionService.SaveTransaction(&mockTransaction)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, response)

		expectedTransaction := &presentation.TransactionDTO{
			TransactionID:   savedTransaction.ID,
			Description:     savedTransaction.Description,
			TransactionDate: util.FormatDate(savedTransaction.TransactionDate),
			PurchaseAmount:  savedTransaction.PurchaseAmount,
		}
		assert.Equal(t, expectedTransaction, response)
	})
	t.Run("Save transaction with success but error on save cache", func(t *testing.T) {
		// given
		mockTransaction := model.Transaction{
			Description:     "mock description",
			TransactionDate: time.Now(),
			PurchaseAmount:  1.0,
		}
		savedTransaction := mockTransaction
		savedTransaction.ID = int64(1)

		// when
		mockRepository.EXPECT().SaveTransaction(&mockTransaction).Return(&savedTransaction, nil)
		mockCache.EXPECT().Save(savedTransaction.ID, &savedTransaction).Return(errors.New("mock error"))

		response, err := transactionService.SaveTransaction(&mockTransaction)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, response)

		expectedTransaction := &presentation.TransactionDTO{
			TransactionID:   savedTransaction.ID,
			Description:     savedTransaction.Description,
			TransactionDate: util.FormatDate(savedTransaction.TransactionDate),
			PurchaseAmount:  savedTransaction.PurchaseAmount,
		}
		assert.Equal(t, expectedTransaction, response)
	})
	t.Run("Save transaction with error on save repository", func(t *testing.T) {
		// given
		mockTransaction := model.Transaction{
			Description:     "mock description",
			TransactionDate: time.Now(),
			PurchaseAmount:  1.0,
		}
		expectedError := presentation.NewApiError(http.StatusInternalServerError, "error saving transaction")

		// when
		mockRepository.EXPECT().SaveTransaction(&mockTransaction).Return(nil, errors.New("mock error"))

		// then
		defer assertPanicApiErrors(t, expectedError)

		_, _ = transactionService.SaveTransaction(&mockTransaction)
	})
}

func Test_TransactionService_UpdateTransactionByID(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	mockRepository := mock_repository.NewMockTransactionRepository(mockController)
	mockCache := mock_repository.NewMockTransactionCache(mockController)

	transactionService := NewTransactionService(slog.Default(), mockRepository, mockCache)

	t.Run("Update transaction by id with success", func(t *testing.T) {
		// given
		mockTransaction := model.Transaction{
			Description:     "mock description",
			TransactionDate: time.Now(),
			PurchaseAmount:  1.0,
		}
		updatedTransaction := mockTransaction
		updatedTransaction.ID = int64(1)

		// when
		mockRepository.EXPECT().UpdateTransaction(updatedTransaction.ID, &mockTransaction).Return(&updatedTransaction, nil)
		mockCache.EXPECT().Save(updatedTransaction.ID, &updatedTransaction).Return(nil)

		response, err := transactionService.UpdateTransactionByID(updatedTransaction.ID, &mockTransaction)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, response)

		expectedTransaction := &presentation.TransactionDTO{
			TransactionID:   updatedTransaction.ID,
			Description:     updatedTransaction.Description,
			TransactionDate: util.FormatDate(updatedTransaction.TransactionDate),
			PurchaseAmount:  updatedTransaction.PurchaseAmount,
		}
		assert.Equal(t, expectedTransaction, response)
	})
	t.Run("Update transaction by id with success but error on save cache", func(t *testing.T) {
		// given
		mockTransaction := model.Transaction{
			Description:     "mock description",
			TransactionDate: time.Now(),
			PurchaseAmount:  1.0,
		}
		updatedTransaction := mockTransaction
		updatedTransaction.ID = int64(1)

		// when
		mockRepository.EXPECT().UpdateTransaction(updatedTransaction.ID, &mockTransaction).Return(&updatedTransaction, nil)
		mockCache.EXPECT().Save(updatedTransaction.ID, &updatedTransaction).Return(errors.New("mock error"))

		response, err := transactionService.UpdateTransactionByID(updatedTransaction.ID, &mockTransaction)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, response)

		expectedTransaction := &presentation.TransactionDTO{
			TransactionID:   updatedTransaction.ID,
			Description:     updatedTransaction.Description,
			TransactionDate: util.FormatDate(updatedTransaction.TransactionDate),
			PurchaseAmount:  updatedTransaction.PurchaseAmount,
		}
		assert.Equal(t, expectedTransaction, response)
	})
	t.Run("Update transaction by id error transaction not found", func(t *testing.T) {
		// given
		mockTransaction := model.Transaction{
			ID: int64(1),
		}
		expectedError := presentation.NewApiError(http.StatusNotFound, "transaction not found")

		// when
		mockRepository.EXPECT().UpdateTransaction(mockTransaction.ID, &mockTransaction).Return(nil, nil)

		// then
		defer assertPanicApiErrors(t, expectedError)

		_, _ = transactionService.UpdateTransactionByID(mockTransaction.ID, &mockTransaction)
	})
	t.Run("Update transaction by id error updating transaction", func(t *testing.T) {
		// given
		mockTransaction := model.Transaction{
			ID: int64(1),
		}
		expectedError := presentation.NewApiError(http.StatusInternalServerError, "error updating transaction")

		// when
		mockRepository.EXPECT().UpdateTransaction(mockTransaction.ID, &mockTransaction).Return(nil, errors.New("mock error"))

		// then
		defer assertPanicApiErrors(t, expectedError)

		_, _ = transactionService.UpdateTransactionByID(mockTransaction.ID, &mockTransaction)
	})
}

func Test_TransactionService_DeleteTransactionByID(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	mockRepository := mock_repository.NewMockTransactionRepository(mockController)
	mockCache := mock_repository.NewMockTransactionCache(mockController)

	transactionService := NewTransactionService(slog.Default(), mockRepository, mockCache)

	t.Run("Delete transaction by id with success", func(t *testing.T) {
		// given
		mockTransaction := model.Transaction{
			ID: int64(1),
		}

		// when
		mockCache.EXPECT().Get(mockTransaction.ID).Return(nil)
		mockRepository.EXPECT().GetTransaction(mockTransaction.ID).Return(&mockTransaction, nil)
		mockRepository.EXPECT().LogicalDeleteTransaction(mockTransaction.ID).Return(&mockTransaction.ID, nil)
		mockCache.EXPECT().Save(mockTransaction.ID, &mockTransaction).Return(nil)

		err := transactionService.DeleteTransactionByID(mockTransaction.ID)

		// then
		assert.NoError(t, err)
	})
	t.Run("Delete transaction by id with success but error on save cache", func(t *testing.T) {
		// given
		mockTransaction := model.Transaction{
			ID: int64(1),
		}

		// when
		mockCache.EXPECT().Get(mockTransaction.ID).Return(nil)
		mockRepository.EXPECT().GetTransaction(mockTransaction.ID).Return(&mockTransaction, nil)
		mockRepository.EXPECT().LogicalDeleteTransaction(mockTransaction.ID).Return(&mockTransaction.ID, nil)
		mockCache.EXPECT().Save(mockTransaction.ID, &mockTransaction).Return(errors.New("mock error"))

		err := transactionService.DeleteTransactionByID(mockTransaction.ID)

		// then
		assert.NoError(t, err)
	})
	t.Run("Delete transaction by id error transaction already deleted in cache", func(t *testing.T) {
		// given
		mockTransaction := model.Transaction{
			ID:      int64(1),
			Deleted: true,
		}
		expectedError := presentation.NewApiError(http.StatusNotFound, "transaction not found")

		// when
		mockCache.EXPECT().Get(mockTransaction.ID).Return(&mockTransaction)

		// then
		defer assertPanicApiErrors(t, expectedError)

		_ = transactionService.DeleteTransactionByID(mockTransaction.ID)
	})
	t.Run("Delete transaction by id error transaction not found in db", func(t *testing.T) {
		// given
		mockTransaction := model.Transaction{
			ID: int64(1),
		}
		expectedError := presentation.NewApiError(http.StatusNotFound, "transaction not found")

		// when
		mockCache.EXPECT().Get(mockTransaction.ID).Return(nil)
		mockRepository.EXPECT().GetTransaction(mockTransaction.ID).Return(nil, nil)

		// then
		defer assertPanicApiErrors(t, expectedError)

		_ = transactionService.DeleteTransactionByID(mockTransaction.ID)
	})
	t.Run("Delete transaction by id error deleting transaction on get transaction", func(t *testing.T) {
		// given
		mockTransaction := model.Transaction{
			ID: int64(1),
		}
		expectedError := presentation.NewApiError(http.StatusInternalServerError, "error deleting transaction")

		// when
		mockCache.EXPECT().Get(mockTransaction.ID).Return(nil)
		mockRepository.EXPECT().GetTransaction(mockTransaction.ID).Return(nil, errors.New("mock error"))

		// then
		defer assertPanicApiErrors(t, expectedError)

		_ = transactionService.DeleteTransactionByID(mockTransaction.ID)
	})
	t.Run("Delete transaction by id error invalid transaction id", func(t *testing.T) {
		// given
		mockTransaction := model.Transaction{
			ID: int64(0),
		}
		expectedError := presentation.NewApiError(http.StatusBadRequest, "invalid transaction id: 0")

		// then
		defer assertPanicApiErrors(t, expectedError)

		// when
		_ = transactionService.DeleteTransactionByID(mockTransaction.ID)
	})
	t.Run("Delete transaction by id error on logical delete repository", func(t *testing.T) {
		// given
		mockTransaction := model.Transaction{
			ID: int64(1),
		}
		expectedError := presentation.NewApiError(http.StatusInternalServerError, "error deleting transaction")

		mockCache.EXPECT().Get(mockTransaction.ID).Return(nil)
		mockRepository.EXPECT().GetTransaction(mockTransaction.ID).Return(&mockTransaction, nil)
		mockRepository.EXPECT().LogicalDeleteTransaction(mockTransaction.ID).Return(nil, errors.New("mock error"))

		// then
		defer assertPanicApiErrors(t, expectedError)

		// when
		_ = transactionService.DeleteTransactionByID(mockTransaction.ID)
	})
}

func assertPanicApiErrors(t *testing.T, expectedError *presentation.ApiError) {
	if r := recover(); r != nil {
		assert.Equal(t, expectedError, r)
	}
}
