package main

import (
	"net/http"

	"github.com/8run0/otel/backend/pkg/api"
)

func main() {
	svr := api.NewServer()
	http.ListenAndServe(":3333", svr)
}
