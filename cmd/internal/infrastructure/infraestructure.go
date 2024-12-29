package infrastructure

import (
	"log/slog"

	"github.com/gorilla/mux"
)

type Infrastructure struct {
	Log      *slog.Logger
	Router   *Routes
	Database *DB
	Cache    *Cache
}

func InitInfrastructure() *Infrastructure {
	log := slog.Default()

	log.Info("Initializing mux router..")
	router := NewRouter(8080, mux.NewRouter())

	log.Info("Initializing database client..")
	database := NewDBClient()

	log.Info("Initializing cache client..")
	cache := NewCache()

	return &Infrastructure{
		Log:      slog.Default(),
		Router:   router,
		Database: database,
		Cache:    cache,
	}
}
