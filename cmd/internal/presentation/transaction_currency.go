package presentation

type TransactionCurrencyDTO struct {
	TransactionID           int64   `json:"transaction_id"`
	Description             string  `json:"description"`
	TransactionDate         string  `json:"transaction_date"`
	PurchaseAmount          float32 `json:"purchase_amount"`
	ExchangeRate            float32 `json:"exchange_rate"`
	ConvertedPurchaseAmount float32 `json:"converted_purchase_amount"`
}
