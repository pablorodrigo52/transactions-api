package repository

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pablorodrigo52/transaction-api/cmd/internal/model"
	"github.com/stretchr/testify/assert"
)

func Test_GetExchangeRateByCountry_APICall(t *testing.T) {
	tests := []struct {
		name           string
		country        string
		mockResponse   string
		mockStatusCode int
		expectedError  error
		expectedResult *model.TreasuryRatesExchange
	}{
		{
			name:           "Should execute call with success",
			country:        "Brazil",
			mockResponse:   `{"data": [{"record_date": "2024-09-30","country": "Brazil","exchange_rate": "5.434","currency": "Real","effective_date": "2024-09-30"}]}`,
			mockStatusCode: http.StatusOK,
			expectedError:  nil,
			expectedResult: &model.TreasuryRatesExchange{
				Data: []model.Data{
					{
						RecordDate:    "2024-09-30",
						Country:       "Brazil",
						ExchangeRate:  "5.434",
						Currency:      "Real",
						EffectiveDate: "2024-09-30",
					},
				},
				Meta: nil,
			},
		},
		{
			name:           "Should return an error because status code is different from 200",
			country:        "Brazil",
			mockResponse:   `{"error":"not found"}`,
			mockStatusCode: http.StatusNotFound,
			expectedError:  errors.New("treasury api call error [status_code:404]"),
			expectedResult: nil,
		},
		{
			name:           "Should return an error because invalid JSON response",
			country:        "Brazil",
			mockResponse:   `invalid json`,
			mockStatusCode: http.StatusOK,
			expectedError:  errors.New("invalid character 'i' looking for beginning of value"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.mockStatusCode)
					w.Write([]byte(tt.mockResponse))
				}))
			defer mockServer.Close()

			repo := NewTreasuryRepository(
				mockServer.URL,
				"/services/api/fiscal_service/v1/accounting/od/rates_of_exchange",
				20*time.Millisecond,
				slog.Default(),
			)
			result, err := repo.GetExchangeRateByCountry(context.TODO(), tt.country)

			if err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func Test_GetExchangeRateByCountry_Client(t *testing.T) {

	t.Run("GetExchangeRateByCountry error on create client", func(t *testing.T) {
		treasuryRepository := NewTreasuryRepository("http://127.0.0.1", "\u2342", 1*time.Second, slog.Default())

		_, err := treasuryRepository.GetExchangeRateByCountry(context.TODO(), "Brazil")

		assert.Error(t, err)
	})
}
