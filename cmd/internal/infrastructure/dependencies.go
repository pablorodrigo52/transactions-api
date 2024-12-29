package infrastructure

import (
	"github.com/pablorodrigo52/transaction-api/cmd/internal/controller"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/repository"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/service"
)

type Dependencies struct {
	PingController        controller.PingController
	TransactionController controller.TransactionController
}

func InitDependencies(infrastructure *Infrastructure) *Dependencies {

	// repositories
	transactionRepository := repository.NewTransactionRepository(infrastructure.Log, infrastructure.Database.Database)
	transactionCache := repository.NewTransactionCache(infrastructure.Cache.Cache)

	// services
	transactionService := service.NewTransactionService(infrastructure.Log, transactionRepository, transactionCache)

	// controllers
	pingController := controller.NewPingController()
	transactionController := controller.NewTransactionController(infrastructure.Log, transactionService)

	return &Dependencies{
		PingController:        *pingController,
		TransactionController: *transactionController,
	}
}
