package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/pablorodrigo52/transaction-api/cmd/internal/model"
)

type TreasuryRepository interface {
	GetExchangeRateByCountry(ctx context.Context, country string) (*model.TreasuryRatesExchange, error)
}

//go:generate mockgen -source=./treasury_repository.go -destination=./mocks/treasury_repository_mock.go

type TreasuryRepositoryImpl struct {
	domain string
	path   string
	client http.Client
	log    *slog.Logger
}

func NewTreasuryRepository(domain, path string, timeout time.Duration, log *slog.Logger) *TreasuryRepositoryImpl {
	return &TreasuryRepositoryImpl{
		domain: domain,
		path:   path,
		client: http.Client{
			Timeout: timeout,
		},
		log: log,
	}
}

func (r *TreasuryRepositoryImpl) GetExchangeRateByCountry(ctx context.Context, country string) (*model.TreasuryRatesExchange, error) {

	completeUrl := fmt.Sprintf(
		"%s%s?fields=record_date,country,exchange_rate,currency,effective_date&filter=country:eq:%s&sort=-record_date&page[number]=1&page[size]=1&format=json",
		r.domain,
		r.path,
		url.QueryEscape(country),
	)

	r.log.Info("Executing api call to", "url", completeUrl)
	resp, err := r.client.Get(completeUrl)
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
