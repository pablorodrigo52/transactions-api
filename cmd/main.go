package main

import (
	"fmt"
	"net/http"

	"github.com/pablorodrigo52/transaction-api/cmd/internal/infrastructure"
	"github.com/pablorodrigo52/transaction-api/cmd/internal/middleware"
)

func main() {
	config := infrastructure.InitInfrastructure()
	dependencies := infrastructure.InitDependencies(config)

	initMiddlewares(config)
	initHandlers(config, dependencies)

	config.Log.Info(fmt.Sprintf("Starting server on http://localhost:%d", config.Router.Port))
	if r := http.ListenAndServe(fmt.Sprintf(":%d", config.Router.Port), config.Router.MuxRouter); r != nil {
		config.Log.Error("Server failed to start", "error", r)
	}
}

func initMiddlewares(config *infrastructure.Infrastructure) {
	config.Router.MuxRouter.Use(middleware.ErrorHandler)
	config.Router.MuxRouter.Use(middleware.JSONContentTypeMiddleware)
}

func initHandlers(config *infrastructure.Infrastructure, dependencies *infrastructure.Dependencies) {
	// ping handler
	config.Router.MuxRouter.HandleFunc("/ping", dependencies.PingController.Ping).Methods("GET")

	// transaction handlers
	r := config.Router.MuxRouter.PathPrefix("/v1").Subrouter()
	r.HandleFunc("/transaction/{id}", dependencies.TransactionController.GetTransactionByID).Methods("GET")
	r.HandleFunc("/transaction", dependencies.TransactionController.CreateTransaction).Methods("POST")
	r.HandleFunc("/transaction/{id}", dependencies.TransactionController.UpdateTransaction).Methods("PUT")
	r.HandleFunc("/transaction/{id}", dependencies.TransactionController.DeleteTransaction).Methods("DELETE")

}
