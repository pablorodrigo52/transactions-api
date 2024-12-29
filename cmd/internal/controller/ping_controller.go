package controller

import (
	"net/http"
)

type PingController struct {
}

func NewPingController() *PingController {
	return &PingController{}
}

func (p *PingController) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
