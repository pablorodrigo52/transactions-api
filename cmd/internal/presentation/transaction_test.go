package presentation

import (
	"net/http"
	"testing"
	"time"

	"github.com/pablorodrigo52/transaction-api/cmd/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestValidateRequest(t *testing.T) {
	tests := []struct {
		name          string
		dto           TransactionDTO
		expectedError *ApiError
	}{
		{
			name: "Validate Request with success",
			dto: TransactionDTO{
				Description:     "Valid Description",
				TransactionDate: "2018-09-26T10:36:40Z",
				PurchaseAmount:  100.0,
			},
			expectedError: nil,
		},
		{
			name: "Validate Request error, empty description",
			dto: TransactionDTO{
				Description:     "",
				TransactionDate: "2018-09-26T10:36:40Z",
				PurchaseAmount:  100.0,
			},
			expectedError: NewApiError(http.StatusBadRequest, "invalid description, it must be between 1 and 50 characters"),
		},
		{
			name: "Validate Request error, description too long",
			dto: TransactionDTO{
				Description:     "This description is way too long and exceeds the fifty character limit",
				TransactionDate: "2018-09-26T10:36:40Z",
				PurchaseAmount:  100.0,
			},
			expectedError: NewApiError(http.StatusBadRequest, "invalid description, it must be between 1 and 50 characters"),
		},
		{
			name: "Validate Request error, empty transaction date",
			dto: TransactionDTO{
				Description:     "Valid Description",
				TransactionDate: "",
				PurchaseAmount:  100.0,
			},
			expectedError: NewApiError(http.StatusBadRequest, "transaction date must not be empty"),
		},
		{
			name: "Validate Request error, invalid transaction date",
			dto: TransactionDTO{
				Description:     "Valid Description",
				TransactionDate: "invalid-date",
				PurchaseAmount:  100.0,
			},
			expectedError: NewApiError(http.StatusBadRequest, "invalid date format expected 2006-01-02T15:04:05Z07:00"),
		},
		{
			name: "Validate Request error, invalid purchase amount",
			dto: TransactionDTO{
				Description:     "Valid Description",
				TransactionDate: "2018-09-26T10:36:40Z",
				PurchaseAmount:  -10.0,
			},
			expectedError: NewApiError(http.StatusBadRequest, "invalid purchase amount, it must be greater than 0"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer assertPanicErrors(t, tt.expectedError)

			if tt.expectedError == nil {
				assert.NotPanics(t, func() {
					tt.dto.Validate()
				})
			} else {
				assert.Panics(t, func() {
					tt.dto.Validate()
				})
			}
		})
	}
}
func TestToTransaction(t *testing.T) {
	tests := []struct {
		name     string
		dto      TransactionDTO
		expected model.Transaction
	}{
		{
			name: "Convert DTO to Transaction with success",
			dto: TransactionDTO{
				TransactionID:   1,
				Description:     "Valid Description",
				TransactionDate: "2018-09-26T10:36:40Z",
				PurchaseAmount:  100.0,
			},
			expected: model.Transaction{
				ID:              1,
				Description:     "Valid Description",
				TransactionDate: time.Date(2018, 9, 26, 10, 36, 40, 0, time.UTC),
				PurchaseAmount:  100.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction := tt.dto.ToTransaction()

			assert.Equal(t, tt.expected.ID, transaction.ID)
			assert.Equal(t, tt.expected.Description, transaction.Description)
			assert.Equal(t, tt.expected.TransactionDate, transaction.TransactionDate)
			assert.Equal(t, tt.expected.PurchaseAmount, transaction.PurchaseAmount)
		})
	}
}

func assertPanicErrors(t *testing.T, expectedError *ApiError) {
	if r := recover(); r != nil {
		assert.Equal(t, expectedError, r)
	}
}
