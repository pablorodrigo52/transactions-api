package repository

import (
	"errors"
	"testing"
	"time"

	"log/slog"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/model"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/util"
	"github.com/stretchr/testify/assert"
)

func Test_TransactionRepository_GetTransaction(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	logger := slog.Default()
	repository := NewTransactionRepository(logger, db)
	selectQuery := "SELECT id, description, transaction_date, purchase_amount, deleted FROM transactions WHERE id = \\?"

	t.Run("GetTransaction with success", func(t *testing.T) {
		// Given
		transactionID := int64(1)
		expectedTransaction := &model.Transaction{
			ID:              transactionID,
			Description:     "Test Transaction",
			TransactionDate: time.Date(2023, 10, 10, 0, 0, 0, 0, time.UTC),
			PurchaseAmount:  100.0,
			Deleted:         false,
		}

		rows := sqlmock.
			NewRows([]string{"id", "description", "transaction_date", "purchase_amount", "deleted"}).
			AddRow(expectedTransaction.ID, expectedTransaction.Description, expectedTransaction.TransactionDate.Format(time.RFC3339), expectedTransaction.PurchaseAmount, expectedTransaction.Deleted)

		mock.ExpectQuery(selectQuery).
			WithArgs(transactionID).
			WillReturnRows(rows)

		// When
		transaction, err := repository.GetTransaction(transactionID)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, transaction)
		assert.Equal(t, expectedTransaction, transaction)
	})

	t.Run("GetTransaction empty due to transaction not found", func(t *testing.T) {
		// Given
		transactionID := int64(2)

		mock.ExpectQuery(selectQuery).
			WithArgs(transactionID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "description", "transaction_date", "purchase_amount", "deleted"}))

		// When
		transaction, err := repository.GetTransaction(transactionID)

		// Then
		assert.NoError(t, err)
		assert.Nil(t, transaction)
		assert.Nil(t, err)
	})

	t.Run("GetTransaction error due to query error", func(t *testing.T) {
		// Given
		transactionID := int64(3)
		expectedErrorMessage := "query error"
		mock.ExpectQuery(selectQuery).
			WithArgs(transactionID).
			WillReturnError(errors.New(expectedErrorMessage))

		// When
		transaction, err := repository.GetTransaction(transactionID)

		// Then
		assert.Error(t, err)
		assert.Nil(t, transaction)
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("GetTransaction error due to scan error", func(t *testing.T) {
		// Given
		transactionID := int64(1)

		rows := sqlmock.
			NewRows([]string{"id", "description", "transaction_date", "purchase_amount", "deleted"}).
			AddRow(nil, nil, nil, nil, nil)

		mock.ExpectQuery(selectQuery).
			WithArgs(transactionID).
			WillReturnRows(rows)

		// When
		transaction, err := repository.GetTransaction(transactionID)

		// Then
		assert.Error(t, err)
		assert.Nil(t, transaction)
		assert.Equal(t, "sql: Scan error on column index 0, name \"id\": converting NULL to int64 is unsupported", err.Error())
	})

	t.Run("GetTransaction error due parse transaction error", func(t *testing.T) {
		// Given
		transactionID := int64(4)
		transactionDate := "invalid-date"

		rows := sqlmock.NewRows([]string{"id", "description", "transaction_date", "purchase_amount", "deleted"}).
			AddRow(transactionID, "Test Transaction", transactionDate, 100.0, false)

		mock.ExpectQuery(selectQuery).
			WithArgs(transactionID).
			WillReturnRows(rows)

		// When
		transaction, err := repository.GetTransaction(transactionID)

		// Then
		assert.Error(t, err)
		assert.Nil(t, transaction)
		assert.Equal(t, "invalid date format expected 2006-01-02T15:04:05Z07:00", err.Error())
	})
}

func Test_TransactionRepository_SaveTransaction(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	logger := slog.Default()
	repository := NewTransactionRepository(logger, db)
	insertQuery := "INSERT INTO transactions \\(description, transaction_date, purchase_amount\\) VALUES \\(\\?, \\?, \\?\\)"

	t.Run("SaveTransaction with success", func(t *testing.T) {
		// Given
		expectedTransaction := &model.Transaction{
			Description:     "test transaction",
			TransactionDate: time.Date(2023, 10, 10, 0, 0, 0, 0, time.UTC),
			PurchaseAmount:  100.0,
		}

		mock.ExpectExec(insertQuery).
			WithArgs(expectedTransaction.Description, util.FormatDate(expectedTransaction.TransactionDate), expectedTransaction.PurchaseAmount).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// When
		transaction, err := repository.SaveTransaction(expectedTransaction)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, transaction)
		assert.Equal(t, expectedTransaction, transaction)
	})

	t.Run("SaveTransaction error on execute query", func(t *testing.T) {
		// Given
		expectedErrorMessage := "mock error run query"
		transaction := &model.Transaction{
			Description:     "test transaction",
			TransactionDate: time.Date(2023, 10, 10, 0, 0, 0, 0, time.UTC),
			PurchaseAmount:  100.0,
		}

		// When
		mock.ExpectExec(insertQuery).
			WithArgs(transaction.Description, util.FormatDate(transaction.TransactionDate), transaction.PurchaseAmount).
			WillReturnError(errors.New(expectedErrorMessage))

		_, err := repository.SaveTransaction(transaction)

		// Then
		assert.Error(t, err)
		assert.Equal(t, expectedErrorMessage, err.Error())
	})
}

func Test_TransactionRepository_UpdateTransaction(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	logger := slog.Default()
	repository := NewTransactionRepository(logger, db)
	updateQuery := "UPDATE transactions SET description = \\?, transaction_date = \\?, purchase_amount = \\? WHERE id = \\? AND deleted = 0"

	t.Run("UpdateTransaction with success", func(t *testing.T) {
		// Given
		transactionID := int64(1)
		expectedTransaction := &model.Transaction{
			ID:              transactionID,
			Description:     "Updated Transaction",
			TransactionDate: time.Date(2023, 10, 10, 0, 0, 0, 0, time.UTC),
			PurchaseAmount:  150.0,
		}

		mock.ExpectExec(updateQuery).
			WithArgs(expectedTransaction.Description, util.FormatDate(expectedTransaction.TransactionDate), expectedTransaction.PurchaseAmount, transactionID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// When
		transaction, err := repository.UpdateTransaction(transactionID, expectedTransaction)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, transaction)
		assert.Equal(t, expectedTransaction, transaction)
	})

	t.Run("UpdateTransaction transaction not found", func(t *testing.T) {
		// Given
		transactionID := int64(2)
		transaction := &model.Transaction{
			ID:              transactionID,
			Description:     "Updated Transaction",
			TransactionDate: time.Date(2023, 10, 10, 0, 0, 0, 0, time.UTC),
			PurchaseAmount:  150.0,
		}

		mock.ExpectExec(updateQuery).
			WithArgs(transaction.Description, util.FormatDate(transaction.TransactionDate), transaction.PurchaseAmount, transactionID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		// When
		updatedTransaction, err := repository.UpdateTransaction(transactionID, transaction)

		// Then
		assert.NoError(t, err)
		assert.Nil(t, updatedTransaction)
	})

	t.Run("UpdateTransaction error on execute query", func(t *testing.T) {
		// Given
		transactionID := int64(3)
		expectedErrorMessage := "mock error run query"
		transaction := &model.Transaction{
			ID:              transactionID,
			Description:     "Updated Transaction",
			TransactionDate: time.Date(2023, 10, 10, 0, 0, 0, 0, time.UTC),
			PurchaseAmount:  150.0,
		}

		mock.ExpectExec(updateQuery).
			WithArgs(transaction.Description, util.FormatDate(transaction.TransactionDate), transaction.PurchaseAmount, transactionID).
			WillReturnError(errors.New(expectedErrorMessage))

		// When
		updatedTransaction, err := repository.UpdateTransaction(transactionID, transaction)

		// Then
		assert.Error(t, err)
		assert.Nil(t, updatedTransaction)
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("UpdateTransaction error on get rows", func(t *testing.T) {
		// Given
		transactionID := int64(4)
		expectedErrorMessage := "mock error rows affected"
		transaction := &model.Transaction{
			ID:              transactionID,
			Description:     "Updated Transaction",
			TransactionDate: time.Date(2023, 10, 10, 0, 0, 0, 0, time.UTC),
			PurchaseAmount:  150.0,
		}

		mock.ExpectExec(updateQuery).
			WithArgs(transaction.Description, util.FormatDate(transaction.TransactionDate), transaction.PurchaseAmount, transactionID).
			WillReturnResult(sqlmock.NewErrorResult(errors.New(expectedErrorMessage)))

		// When
		updatedTransaction, err := repository.UpdateTransaction(transactionID, transaction)

		// Then
		assert.Error(t, err)
		assert.Nil(t, updatedTransaction)
		assert.Equal(t, expectedErrorMessage, err.Error())
	})
}

func Test_TransactionRepository_LogicalDeleteTransaction(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	logger := slog.Default()
	repository := NewTransactionRepository(logger, db)
	deleteQuery := "UPDATE transactions SET deleted = 1 WHERE id = \\? AND deleted = 0"

	t.Run("LogicalDeleteTransaction with success", func(t *testing.T) {
		// Given
		transactionID := int64(1)

		mock.ExpectExec(deleteQuery).
			WithArgs(transactionID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// When
		deletedID, err := repository.LogicalDeleteTransaction(transactionID)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, deletedID)
		assert.Equal(t, transactionID, *deletedID)
	})

	t.Run("LogicalDeleteTransaction transaction not found", func(t *testing.T) {
		// Given
		transactionID := int64(2)

		mock.ExpectExec(deleteQuery).
			WithArgs(transactionID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		// When
		deletedID, err := repository.LogicalDeleteTransaction(transactionID)

		// Then
		assert.NoError(t, err)
		assert.Nil(t, deletedID)
	})

	t.Run("LogicalDeleteTransaction error on execute query", func(t *testing.T) {
		// Given
		transactionID := int64(3)
		expectedErrorMessage := "mock error run query"

		mock.ExpectExec(deleteQuery).
			WithArgs(transactionID).
			WillReturnError(errors.New(expectedErrorMessage))

		// When
		deletedID, err := repository.LogicalDeleteTransaction(transactionID)

		// Then
		assert.Error(t, err)
		assert.Nil(t, deletedID)
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("LogicalDeleteTransaction error on get rows affected", func(t *testing.T) {
		// Given
		transactionID := int64(4)
		expectedErrorMessage := "mock error rows affected"

		mock.ExpectExec(deleteQuery).
			WithArgs(transactionID).
			WillReturnResult(sqlmock.NewErrorResult(errors.New(expectedErrorMessage)))

		// When
		deletedID, err := repository.LogicalDeleteTransaction(transactionID)

		// Then
		assert.Error(t, err)
		assert.Nil(t, deletedID)
		assert.Equal(t, expectedErrorMessage, err.Error())
	})
}
