package service

import (
	"log/slog"

	"github.com/pablorodrigo52/transaction-api/cmd/internal/model"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/presentation"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/repository"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/util"
)

type TransactionService interface {
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
