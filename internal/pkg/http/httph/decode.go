package httph

import (
	"encoding/json"
	"net/http"
)

func DecodeJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}
