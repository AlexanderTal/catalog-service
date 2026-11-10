package httph

import (
	"context"
	"net/http"
)

type contextKeyError struct{}

type contextValueError struct {
	err        error
	statusCode int
}

func errorPrepare(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyError{}, &contextValueError{})
}

func errorApply(ctx context.Context, err error) {
	// Достаем указатель из контекста и проверяем, что он там есть
	if v, ok := ctx.Value(contextKeyError{}).(*contextValueError); ok && v != nil {
		v.err = err
	}
}

func errorGet(ctx context.Context) error {
	if v, ok := ctx.Value(contextKeyError{}).(*contextValueError); ok && v != nil {
		return v.err
	}
	return nil
}

func errorApplyStatusCode(ctx context.Context, statusCode int) {
	if v, ok := ctx.Value(contextKeyError{}).(*contextValueError); ok && v != nil {
		v.statusCode = statusCode
	}
}

func errorGetStatusCode(ctx context.Context) int {
	if v, ok := ctx.Value(contextKeyError{}).(*contextValueError); ok && v != nil {
		return v.statusCode
	}
	return 0
}

// Публичные обёртки (из подсказки к заданию)
type Middleware = func(http.Handler) http.Handler

func ErrorPrepare(r *http.Request) *http.Request {
	return r.WithContext(errorPrepare(r.Context()))
}

func ErrorGet(r *http.Request) error {
	return errorGet(r.Context())
}

func ErrorApply(r *http.Request, err error) {
	errorApply(r.Context(), err)
}

func ErrorApplyStatusCode(r *http.Request, statusCode int) {
	errorApplyStatusCode(r.Context(), statusCode)
}

func ErrorGetStatusCode(r *http.Request) int {
	return errorGetStatusCode(r.Context())
}

func NewErrorMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, ErrorPrepare(r))
		})
	}
}
