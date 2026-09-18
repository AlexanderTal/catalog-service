package rhandler

import "net/http"

type Health interface {
	LastCheck(w http.ResponseWriter, r *http.Request)
}
