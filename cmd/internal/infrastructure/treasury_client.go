package infrastructure

import "time"

type TreasuryClient struct {
	domain  string
	path    string
	timeout time.Duration
}

func NewTreasuryClient() *TreasuryClient {
	return &TreasuryClient{
		domain:  "https://api.fiscaldata.treasury.gov",
		path:    "/services/api/fiscal_service/v1/accounting/od/rates_of_exchange",
		timeout: 30 * time.Second,
	}
}
