package repository

import (
	"database/sql"
	"log/slog"

	"github.com/pablorodrigo52/transaction-api/cmd/internal/model"
)

type TransactionRepository interface {
	SaveTransaction(transaction *model.Transaction) (*model.Transaction, error)
}

type TransactionRepositoryImpl struct {
	log *slog.Logger
	db  *sql.DB
}

func NewTransactionRepository(log *slog.Logger, db *sql.DB) *TransactionRepositoryImpl {
	return &TransactionRepositoryImpl{
		log: log,
		db:  db,
	}
}

func (t *TransactionRepositoryImpl) SaveTransaction(transaction *model.Transaction) (*model.Transaction, error) {
	trx, err := t.db.Exec("INSERT INTO transactions (description, transaction_date, purchase_amount) VALUES (?, ?, ?)", transaction.Description, transaction.TransactionDate, transaction.PurchaseAmount)
	if err != nil {
		return nil, err
	}

	transaction.ID, _ = trx.LastInsertId()
	return transaction, nil
}
