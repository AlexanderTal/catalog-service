package util

import (
	"net/http"
	"strings"
)

// IsFilteredHttpRoute возвращает true для маршрутов, которые не нужно логировать
func IsFilteredHttpRoute(r *http.Request) bool {
	return strings.Contains(r.RequestURI, "health") ||
		strings.Contains(r.RequestURI, "debug") ||
		strings.Contains(r.RequestURI, "metric")
}
