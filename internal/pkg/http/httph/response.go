package httph

import (
	"encoding/json"
	"errors"
	"net/http"
)

type httpCoder interface {
	error
	HTTPStatus() int
}

func SendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func SendEmpty(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func sendError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

func HandleError(w http.ResponseWriter, err error) {
	var hc httpCoder
	if errors.As(err, &hc) {
		sendError(w, hc.HTTPStatus(), hc)
		return
	}
	sendError(w, http.StatusInternalServerError, err)
}
