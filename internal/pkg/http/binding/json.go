package binding

import (
	"errors"
	"net/http"

	"github.com/AlexanderTal/catalog-service/internal/pkg/http/httph"
)

type jsonBinding struct{}

func (jsonBinding) Name() string {
	return "JSON"
}

func (jsonBinding) Bind(req *http.Request, obj any) error {
	if req == nil || req.Body == nil {
		return errors.New("invalid request")
	}

	if err := httph.DecodeJSON(req, obj); err != nil {
		return err
	}

	return validate(obj)
}
