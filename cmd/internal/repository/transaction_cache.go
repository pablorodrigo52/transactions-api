package repository

import (
	"errors"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/model"
)

type TransactionCache interface {
	Get(transactionID int64) *model.Transaction
	Save(transactionID int64, transaction *model.Transaction) error
}

//go:generate mockgen -source=./transaction_cache.go -destination=./mocks/transaction_cache_mock.go

type TransactionCacheImpl struct {
	cache *ristretto.Cache
	TTL   time.Duration
	Cost  int64
}

func NewTransactionCache(cache *ristretto.Cache) *TransactionCacheImpl {
	return &TransactionCacheImpl{
		cache: cache,
		TTL:   1 * time.Hour,
		Cost:  1,
	}
}

func (t *TransactionCacheImpl) Get(transactionID int64) *model.Transaction {
	if transaction, found := t.cache.Get(transactionID); found {
		return transaction.(*model.Transaction)
	}

	return nil
}

func (t *TransactionCacheImpl) Save(transactionID int64, transaction *model.Transaction) error {
	t.cache.Del(transactionID)
	if !t.cache.SetWithTTL(transactionID, transaction, t.Cost, t.TTL) {
		errorMessage := "error saving transaction in cache"
		return errors.New(errorMessage)
	}

	return nil
}
