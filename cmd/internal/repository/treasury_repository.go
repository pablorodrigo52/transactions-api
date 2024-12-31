package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/pablorodrigo52/transaction-api/cmd/internal/model"
)

type TreasuryRepository interface {
	GetExchangeRateByCountry(country string) (*model.TreasuryRatesExchange, error)
}

//go:generate mockgen -source=./treasury_repository.go -destination=./mocks/treasury_repository_mock.go

type TreasuryRepositoryImpl struct {
	Domain string
	Path   string
	Client http.Client
	log    *slog.Logger
}

func NewTreasuryRepository(domain, path string, timeout time.Duration, log *slog.Logger) *TreasuryRepositoryImpl {
	return &TreasuryRepositoryImpl{
		Domain: domain,
		Path:   path,
		Client: http.Client{
			Timeout: timeout,
		},
		log: log,
	}
}

func (r *TreasuryRepositoryImpl) GetExchangeRateByCountry(country string) (*model.TreasuryRatesExchange, error) {

	url := fmt.Sprintf(
		"%s%s?fields=record_date,country,exchange_rate,currency,effective_date&filter=country:eq:%s&sort=-record_date&page[number]=1&page[size]=1&format=json",
		r.Domain,
		r.Path,
		country,
	)

	resp, err := r.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		r.log.Error("Error on execute treasury api call", "status_code", resp.StatusCode, "response", resp.Body)
		return nil, fmt.Errorf("treasury api call error [status_code:%d]", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var exchangeRateResponse model.TreasuryRatesExchange
	if err := json.Unmarshal(body, &exchangeRateResponse); err != nil {
		return nil, err
	}

	return &exchangeRateResponse, nil
}
