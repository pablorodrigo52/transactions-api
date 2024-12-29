package repository

import (
	"database/sql"
	"log/slog"

	"github.com/pablorodrigo52/transaction-api/cmd/internal/model"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/util"
)

type TransactionRepository interface {
	GetTransaction(transactionID int64) (*model.Transaction, error)
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

func (t *TransactionRepositoryImpl) GetTransaction(transactionID int64) (*model.Transaction, error) {
	result, err := t.db.Query("SELECT id, description, transaction_date, purchase_amount FROM transactions WHERE id = ?", transactionID)

	if err != nil {
		return nil, err
	}

	var transaction model.Transaction
	if result.Next() {
		var transactionDate string
		err := result.Scan(&transaction.ID, &transaction.Description, &transactionDate, &transaction.PurchaseAmount)
		if err != nil {
			return nil, err
		}

		transaction.TransactionDate, err = util.ParseDate(transactionDate)
		if err != nil {
			return nil, err
		}
	}

	return &transaction, nil
}

func (t *TransactionRepositoryImpl) SaveTransaction(transaction *model.Transaction) (*model.Transaction, error) {
	trx, err := t.db.Exec("INSERT INTO transactions (description, transaction_date, purchase_amount) VALUES (?, ?, ?)", transaction.Description, transaction.TransactionDate, transaction.PurchaseAmount)
	if err != nil {
		return nil, err
	}

	transaction.ID, _ = trx.LastInsertId()
	return transaction, nil
}
