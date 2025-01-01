package presentation

import (
	"net/http"
	"strconv"
)

type TransactionID string

func (t *TransactionID) Validate() {
	if t == nil || *t == "" {
		panic(NewApiError(http.StatusBadRequest, "transaction ID is required"))
	}

	id, err := strconv.ParseInt(string(*t), 10, 64)
	if err != nil {
		panic(NewApiError(http.StatusBadRequest, "transaction ID must be a valid number: "+err.Error()))
	}

	if id <= 0 {
		panic(NewApiError(http.StatusBadRequest, "transaction ID must be a valid number"))
	}
}

func (t *TransactionID) Get() int64 {
	id, _ := strconv.ParseInt(string(*t), 10, 64)
	return id
}
