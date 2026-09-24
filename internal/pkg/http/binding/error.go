package binding

import "net/http"

type bindingError struct {
	msg string
}

func (e *bindingError) Error() string   { return e.msg }
func (e *bindingError) HTTPStatus() int { return http.StatusBadRequest }
