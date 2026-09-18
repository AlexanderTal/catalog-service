package rprocessor

import (
	"net/http"

	"github.com/gorilla/mux"
)

func reg(r *mux.Router, method, path string, handler http.Handler) {
	r.Methods(method).Path(path).Handler(handler)
}
