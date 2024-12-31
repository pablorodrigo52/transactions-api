package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/presentation"
	mock_service "github.com/pablorodrigo52/transaction-api/cmd/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateTransactionID(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	mockService := mock_service.NewMockTransactionService(mockController)

	logger := slog.Default()
	controller := NewTransactionController(logger, mockService)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/transactions/{id}", controller.GetTransactionByID)

	t.Run("Validate transaction id with success", func(t *testing.T) {
		// Given
		req, err := http.NewRequest("GET", "/transactions/1", nil)
		assert.NoError(t, err)

		expectedResponse := presentation.TransactionDTO{TransactionID: 1}

		mockService.EXPECT().GetTransactionByID(int64(1)).Return(&expectedResponse, nil)

		// When
		router.ServeHTTP(rr, req)

		// Then
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotNil(t, rr.Body.String())

		var response presentation.TransactionDTO
		err = json.Unmarshal(rr.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
	})

	t.Run("Validate transaction id with error invalid parameter non-integer", func(t *testing.T) {
		// Given
		req, err := http.NewRequest("GET", "/transactions/mock", nil)
		assert.NoError(t, err)

		expectedError := presentation.NewApiError(http.StatusBadRequest, "Transaction ID must be a valid number: strconv.ParseInt: parsing \"mock\": invalid syntax")

		// Then
		defer assertPanicErrors(t, expectedError)

		// When
		router.ServeHTTP(rr, req)
	})

	t.Run("Validate transaction id with error invalid parameter integer inalid", func(t *testing.T) {
		// Given
		req, err := http.NewRequest("GET", "/transactions/0", nil)
		assert.NoError(t, err)

		expectedError := presentation.NewApiError(http.StatusBadRequest, "Transaction ID must be a valid number")

		// Then
		defer assertPanicErrors(t, expectedError)

		// When
		router.ServeHTTP(rr, req)
	})
}

func Test_DecodeTransactionDTO(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	mockService := mock_service.NewMockTransactionService(mockController)

	logger := slog.Default()
	controller := NewTransactionController(logger, mockService)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		controller.decodeTransactionDTO(r)
	})

	t.Run("Decode transaction DTO with success", func(t *testing.T) {
		// Given
		transactionDTO := presentation.TransactionDTO{
			TransactionID:   1,
			Description:     "mock",
			TransactionDate: "2018-09-26T10:36:40Z",
			PurchaseAmount:  1.0,
		}

		body, err := json.Marshal(transactionDTO)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
		assert.NoError(t, err)

		// When
		router.ServeHTTP(rr, req)

		// Then
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Decode transaction DTO with error invalid JSON", func(t *testing.T) {
		// Given
		req, err := http.NewRequest("POST", "/transactions", bytes.NewBuffer([]byte("{invalid json}")))
		assert.NoError(t, err)

		expectedError := presentation.NewApiError(http.StatusBadRequest, "Error decoding request body: invalid character 'i' looking for beginning of object key string")

		// Then
		defer assertPanicErrors(t, expectedError)

		// When
		router.ServeHTTP(rr, req)
	})

	t.Run("Decode transaction DTO with error invalid request body", func(t *testing.T) {
		// Given
		transactionDTO := presentation.TransactionDTO{}
		body, err := json.Marshal(transactionDTO)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
		assert.NoError(t, err)

		expectedError := presentation.NewApiError(http.StatusBadRequest, "Error validating request body: invalid description, it must be between 1 and 50 characters")

		// Then
		defer assertPanicErrors(t, expectedError)

		// When
		router.ServeHTTP(rr, req)
	})
}

func Test_GetTransactionByID(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	mockService := mock_service.NewMockTransactionService(mockController)

	logger := slog.Default()
	controller := NewTransactionController(logger, mockService)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/transactions/{id}", controller.GetTransactionByID)

	t.Run("Get transaction by id with success", func(t *testing.T) {
		// Given
		req, err := http.NewRequest("GET", "/transactions/1", nil)
		assert.NoError(t, err)

		expectedResponse := presentation.TransactionDTO{TransactionID: 1}

		mockService.EXPECT().GetTransactionByID(int64(1)).Return(&expectedResponse, nil)

		// When
		router.ServeHTTP(rr, req)

		// Then
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotNil(t, rr.Body.String())

		var response presentation.TransactionDTO
		err = json.Unmarshal(rr.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
	})

	t.Run("Get transaction by id error getting transaction", func(t *testing.T) {
		// Given
		req, err := http.NewRequest("GET", "/transactions/1", nil)
		assert.NoError(t, err)

		expectedError := presentation.NewApiError(http.StatusInternalServerError, "Error getting transaction: mock error")

		mockService.EXPECT().GetTransactionByID(int64(1)).Return(nil, errors.New("mock error"))

		// Then
		defer assertPanicErrors(t, expectedError)

		// When
		router.ServeHTTP(rr, req)
	})
}

func Test_CreateTransaction(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	mockService := mock_service.NewMockTransactionService(mockController)

	logger := slog.Default()
	controller := NewTransactionController(logger, mockService)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/transactions", controller.CreateTransaction).Methods("POST")

	t.Run("Create transaction with success", func(t *testing.T) {
		// Given
		expectedResponse := presentation.TransactionDTO{
			TransactionID:   1,
			Description:     "mock",
			TransactionDate: "2018-09-26T10:36:40Z",
			PurchaseAmount:  1.0,
		}

		body, err := json.Marshal(expectedResponse)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
		assert.NoError(t, err)

		mockService.EXPECT().SaveTransaction(expectedResponse.ToTransaction()).Return(&expectedResponse, nil)

		// When
		router.ServeHTTP(rr, req)

		// Then
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotNil(t, rr.Body.String())

		var response presentation.TransactionDTO
		err = json.Unmarshal(rr.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
	})

	t.Run("Create transaction with error saving transaction", func(t *testing.T) {
		// Given
		transactionDTO := presentation.TransactionDTO{
			TransactionID:   1,
			Description:     "mock",
			TransactionDate: "2018-09-26T10:36:40Z",
			PurchaseAmount:  1.0,
		}

		body, err := json.Marshal(transactionDTO)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
		assert.NoError(t, err)

		expectedError := presentation.NewApiError(http.StatusInternalServerError, "Error saving transaction: mock error")

		mockService.EXPECT().SaveTransaction(transactionDTO.ToTransaction()).Return(nil, errors.New("mock error"))

		// Then
		defer assertPanicErrors(t, expectedError)

		// When
		router.ServeHTTP(rr, req)
	})
}

func Test_UpdateTransaction(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	mockService := mock_service.NewMockTransactionService(mockController)

	logger := slog.Default()
	controller := NewTransactionController(logger, mockService)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/transactions/{id}", controller.UpdateTransaction).Methods("PUT")

	t.Run("Update transaction with success", func(t *testing.T) {
		// Given
		transactionDTO := presentation.TransactionDTO{
			TransactionID:   1,
			Description:     "updated description",
			TransactionDate: "2018-09-26T10:36:40Z",
			PurchaseAmount:  2.0,
		}

		body, err := json.Marshal(transactionDTO)
		assert.NoError(t, err)

		req, err := http.NewRequest("PUT", "/transactions/1", bytes.NewBuffer(body))
		assert.NoError(t, err)

		mockService.EXPECT().UpdateTransactionByID(int64(1), transactionDTO.ToTransaction()).Return(&transactionDTO, nil)

		// When
		router.ServeHTTP(rr, req)

		// Then
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotNil(t, rr.Body.String())

		var response presentation.TransactionDTO
		err = json.Unmarshal(rr.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, transactionDTO, response)
	})

	t.Run("Update transaction with error updating transaction", func(t *testing.T) {
		// Given
		transactionDTO := presentation.TransactionDTO{
			TransactionID:   1,
			Description:     "updated description",
			TransactionDate: "2018-09-26T10:36:40Z",
			PurchaseAmount:  2.0,
		}

		body, err := json.Marshal(transactionDTO)
		assert.NoError(t, err)

		req, err := http.NewRequest("PUT", "/transactions/1", bytes.NewBuffer(body))
		assert.NoError(t, err)

		expectedError := presentation.NewApiError(http.StatusInternalServerError, "Error updating transaction: mock error")

		mockService.EXPECT().UpdateTransactionByID(int64(1), transactionDTO.ToTransaction()).Return(nil, errors.New("mock error"))

		// Then
		defer assertPanicErrors(t, expectedError)

		// When
		router.ServeHTTP(rr, req)
	})
}

func Test_DeleteTransaction(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	mockService := mock_service.NewMockTransactionService(mockController)

	logger := slog.Default()
	controller := NewTransactionController(logger, mockService)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/transactions/{id}", controller.DeleteTransaction).Methods("DELETE")

	t.Run("Delete transaction with success", func(t *testing.T) {
		// Given
		req, err := http.NewRequest("DELETE", "/transactions/1", nil)
		assert.NoError(t, err)

		mockService.EXPECT().DeleteTransactionByID(int64(1)).Return(nil)

		// When
		router.ServeHTTP(rr, req)

		// Then
		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("Delete transaction with error deleting transaction", func(t *testing.T) {
		// Given
		req, err := http.NewRequest("DELETE", "/transactions/1", nil)
		assert.NoError(t, err)

		expectedError := presentation.NewApiError(http.StatusInternalServerError, "Error deleting transaction: mock error")

		mockService.EXPECT().DeleteTransactionByID(int64(1)).Return(errors.New("mock error"))

		// Then
		defer assertPanicErrors(t, expectedError)

		// When
		router.ServeHTTP(rr, req)
	})
}

func assertPanicErrors(t *testing.T, expectedError *presentation.ApiError) {
	if r := recover(); r != nil {
		assert.Equal(t, expectedError, r)
	}
}
