package infrastructure

import (
	"github.com/pablorodrigo52/transaction-api/cmd/internal/controller"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/repository"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/service"
)

type Dependencies struct {
	PingController                controller.PingController
	TransactionController         controller.TransactionController
	TransactionCurrencyController controller.TransactionCurrencyController
}

func InitDependencies(infrastructure *Infrastructure) *Dependencies {

	// repositories
	transactionRepository := repository.NewTransactionRepository(infrastructure.Log, infrastructure.Database.Database)
	transactionCache := repository.NewTransactionCache(infrastructure.Cache.Cache)
	treasuryRepository := repository.NewTreasuryRepository(
		infrastructure.TreasuryClient.domain,
		infrastructure.TreasuryClient.path,
		infrastructure.TreasuryClient.timeout,
		infrastructure.Log)

	// services
	transactionService := service.NewTransactionService(infrastructure.Log, transactionRepository, transactionCache)
	transactionCurrencyService := service.NewTransactionCurrencyService(treasuryRepository, transactionRepository, infrastructure.Log)

	// controllers
	pingController := controller.NewPingController()
	transactionController := controller.NewTransactionController(infrastructure.Log, transactionService)
	transactionCurrencyController := controller.NewTransactionCurrencyController(transactionCurrencyService, infrastructure.Log)

	return &Dependencies{
		PingController:                *pingController,
		TransactionController:         *transactionController,
		TransactionCurrencyController: *transactionCurrencyController,
	}
}
