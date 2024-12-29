package service

import (
	"fmt"
	"log/slog"

	"github.com/pablorodrigo52/transaction-api/cmd/internal/model"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/presentation"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/repository"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/util"
)

type TransactionService interface {
	GetTransactionByID(transactionID int64) (*presentation.TransactionDTO, error)
	SaveTransaction(transaction *model.Transaction) (*presentation.TransactionDTO, error)
}

type TransactionServiceImpl struct {
	log        *slog.Logger
	repository repository.TransactionRepository
	cache      repository.TransactionCache
}

func NewTransactionService(
	log *slog.Logger,
	repository repository.TransactionRepository,
	cache repository.TransactionCache) *TransactionServiceImpl {

	return &TransactionServiceImpl{
		log:        log,
		repository: repository,
		cache:      cache,
	}
}

func (t *TransactionServiceImpl) GetTransactionByID(transactionID int64) (*presentation.TransactionDTO, error) {

	if transactionID <= 0 {
		return nil, fmt.Errorf("invalid transaction id: %d", transactionID)
	}

	// recover from cache
	if trx := t.cache.Get(transactionID); trx != nil {
		t.log.Info("Transaction found in cache", "transaction_id", transactionID)
		return &presentation.TransactionDTO{
			TransactionID:   trx.ID,
			Description:     trx.Description,
			TransactionDate: util.FormatDate(trx.TransactionDate),
			PurchaseAmount:  trx.PurchaseAmount,
		}, nil
	}

	// if not found on cache, go to database
	t.log.Info("Transaction not found in cache, searching on db", "transaction_id", transactionID)
	trx, err := t.repository.GetTransaction(transactionID)
	if err != nil {
		t.log.Error("error getting transaction", "error", err)
		return nil, err
	}

	// if find on database, save on cache
	if err := t.cache.Save(trx.ID, trx); err != nil {
		t.log.Error("error saving transaction cache ", "transaction_id", trx.ID)
	}

	return &presentation.TransactionDTO{
		TransactionID:   trx.ID,
		Description:     trx.Description,
		TransactionDate: util.FormatDate(trx.TransactionDate),
		PurchaseAmount:  trx.PurchaseAmount,
	}, nil
}

func (t *TransactionServiceImpl) SaveTransaction(transaction *model.Transaction) (*presentation.TransactionDTO, error) {

	trx, err := t.repository.SaveTransaction(transaction)
	if err != nil {
		t.log.Error("error saving transaction", "error", err)
		return nil, err
	}

	if err := t.cache.Save(trx.ID, trx); err != nil {
		t.log.Error("error saving transaction cache ", "transaction_id", trx.ID)
	}

	t.log.Info("Transaction saved", "transaction_id", trx.ID)

	return &presentation.TransactionDTO{
		TransactionID:   trx.ID,
		Description:     trx.Description,
		TransactionDate: util.FormatDate(trx.TransactionDate),
		PurchaseAmount:  trx.PurchaseAmount,
	}, nil
}
