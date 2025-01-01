package model

type TreasuryRatesExchange struct {
	Data  []Data             `json:"data"`
	Meta  *Meta              `json:"meta"`
	Links map[string]*string `json:"links"`
}

type Data struct {
	RecordDate    string `json:"record_date"`
	Country       string `json:"country"`
	ExchangeRate  string `json:"exchange_rate"`
	Currency      string `json:"currency"`
	EffectiveDate string `json:"effective_date"`
}

type Meta struct {
	Count       int               `json:"count"`
	DataFormats map[string]string `json:"dataFormats"`
	TotalCount  int               `json:"total-count"`
	TotalPages  int               `json:"total-pages"`
}
