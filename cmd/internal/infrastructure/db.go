package infrastructure

import (
	"database/sql"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/presentation"
)

type DB struct {
	Database *sql.DB
}

func NewDBClient() *DB {
	db, err := sql.Open("sqlite3", "../db/transactions.db")
	if err != nil {
		panic(presentation.NewApiError(http.StatusInternalServerError, "Failed to open database: "+err.Error()))
	}

	initScript, err := os.ReadFile("../scripts/init.sql")
	if err != nil {
		panic(presentation.NewApiError(http.StatusInternalServerError, "Failed to read init.sql: "+err.Error()))
	}

	if _, err := db.Exec(string(initScript)); err != nil {
		panic(presentation.NewApiError(http.StatusInternalServerError, "Failed to execute init.sql: "+err.Error()))
	}

	return &DB{
		Database: db,
	}
}
