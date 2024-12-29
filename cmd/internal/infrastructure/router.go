package infrastructure

import "github.com/gorilla/mux"

type Routes struct {
	Port   int
	MuxRouter *mux.Router
}

func NewRouter(port int, router *mux.Router) *Routes {
	return &Routes{
		Port:   port,
		MuxRouter: router,
	}
}
