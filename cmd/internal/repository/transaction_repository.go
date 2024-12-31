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
	UpdateTransaction(transactionID int64, transaction *model.Transaction) (*model.Transaction, error)
	LogicalDeleteTransaction(transactionID int64) (*int64, error)
}

//go:generate mockgen -source=./transaction_repository.go -destination=./mocks/transaction_repository_mock.go

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
	result, err := t.db.Query("SELECT id, description, transaction_date, purchase_amount, deleted FROM transactions WHERE id = ?", transactionID)
	if err != nil {
		return nil, err
	}

	defer result.Close()

	if result.Next() {
		var transaction model.Transaction
		var transactionDate string

		err := result.Scan(&transaction.ID, &transaction.Description, &transactionDate, &transaction.PurchaseAmount, &transaction.Deleted)
		if err != nil {
			return nil, err
		}

		transaction.TransactionDate, err = util.ParseDate(transactionDate)
		if err != nil {
			return nil, err
		}

		return &transaction, nil
	}

	return nil, nil
}

func (t *TransactionRepositoryImpl) SaveTransaction(transaction *model.Transaction) (*model.Transaction, error) {
	trx, err := t.db.Exec("INSERT INTO transactions (description, transaction_date, purchase_amount) VALUES (?, ?, ?)", transaction.Description, util.FormatDate(transaction.TransactionDate), transaction.PurchaseAmount)
	if err != nil {
		return nil, err
	}

	transaction.ID, _ = trx.LastInsertId()
	return transaction, nil
}

func (t *TransactionRepositoryImpl) UpdateTransaction(transactionID int64, transaction *model.Transaction) (*model.Transaction, error) {
	trx, err := t.db.Exec("UPDATE transactions SET description = ?, transaction_date = ?, purchase_amount = ? WHERE id = ? AND deleted = 0", transaction.Description, util.FormatDate(transaction.TransactionDate), transaction.PurchaseAmount, transactionID)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := trx.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, nil
	}

	return transaction, nil
}

func (t *TransactionRepositoryImpl) LogicalDeleteTransaction(transactionID int64) (*int64, error) {
	trx, err := t.db.Exec("UPDATE transactions SET deleted = 1 WHERE id = ? AND deleted = 0", transactionID)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := trx.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, nil
	}

	return &transactionID, nil
}
