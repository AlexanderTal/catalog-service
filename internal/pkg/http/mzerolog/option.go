package mzerolog

import (
	"net/http"

	"github.com/rs/zerolog"
)

type Option = func(m *middleware)

func WithLogger(l zerolog.Logger) Option {
	return func(m *middleware) {
		m.log = l
	}
}

func WithSkipper(skipper func(r *http.Request) bool) Option {
	return func(m *middleware) {
		if skipper != nil {
			m.fromOptions.skipper = skipper
		}
	}
}
