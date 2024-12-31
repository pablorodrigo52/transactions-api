package presentation

import (
	"testing"
	"time"

	"github.com/pablorodrigo52/transaction-api/cmd/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestValidateRequest(t *testing.T) {
	tests := []struct {
		name          string
		dto           TransactionDTO
		expectedError string
	}{
		{
			name: "Validate Request with success",
			dto: TransactionDTO{
				Description:     "Valid Description",
				TransactionDate: "2018-09-26T10:36:40Z",
				PurchaseAmount:  100.0,
			},
			expectedError: "",
		},
		{
			name: "Validate Request error, empty description",
			dto: TransactionDTO{
				Description:     "",
				TransactionDate: "2018-09-26T10:36:40Z",
				PurchaseAmount:  100.0,
			},
			expectedError: "invalid description, it must be between 1 and 50 characters",
		},
		{
			name: "Validate Request error, description too long",
			dto: TransactionDTO{
				Description:     "This description is way too long and exceeds the fifty character limit",
				TransactionDate: "2018-09-26T10:36:40Z",
				PurchaseAmount:  100.0,
			},
			expectedError: "invalid description, it must be between 1 and 50 characters",
		},
		{
			name: "Validate Request error, empty transaction date",
			dto: TransactionDTO{
				Description:     "Valid Description",
				TransactionDate: "",
				PurchaseAmount:  100.0,
			},
			expectedError: "transaction date must not be empty",
		},
		{
			name: "Validate Request error, invalid transaction date",
			dto: TransactionDTO{
				Description:     "Valid Description",
				TransactionDate: "invalid-date",
				PurchaseAmount:  100.0,
			},
			expectedError: "invalid date format expected 2006-01-02T15:04:05Z07:00",
		},
		{
			name: "Validate Request error, invalid purchase amount",
			dto: TransactionDTO{
				Description:     "Valid Description",
				TransactionDate: "2018-09-26T10:36:40Z",
				PurchaseAmount:  -10.0,
			},
			expectedError: "invalid purchase amount, it must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dto.ValidateRequest()

			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
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

func TestRoundPurchaseAmount(t *testing.T) {
	tests := []struct {
		name     string
		dto      TransactionDTO
		expected float32
	}{
		{
			name: "Round purchase amount with no rounding needed",
			dto: TransactionDTO{
				PurchaseAmount: 100.0,
			},
			expected: 100.0,
		},
		{
			name: "Round purchase amount with rounding down",
			dto: TransactionDTO{
				PurchaseAmount: 100.123,
			},
			expected: 100.12,
		},
		{
			name: "Round purchase amount with rounding up",
			dto: TransactionDTO{
				PurchaseAmount: 100.126,
			},
			expected: 100.13,
		},
		{
			name: "Round purchase amount with exact half",
			dto: TransactionDTO{
				PurchaseAmount: 100.125,
			},
			expected: 100.13,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roundedAmount := tt.dto.RoundPurchaseAmount()
			assert.Equal(t, tt.expected, roundedAmount)
		})
	}
}
