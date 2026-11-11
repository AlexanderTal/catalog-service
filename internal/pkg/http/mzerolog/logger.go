package mzerolog

import (
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/AlexanderTal/catalog-service/internal/pkg/http/httph"
)

type middleware struct {
	log zerolog.Logger

	fromOptions struct {
		skipper func(r *http.Request) bool
	}
}

func (m *middleware) Callback(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const (
			tailSuccess = " finished with no error"
			tailFail    = " finished (or aborted) with error"
		)

		start := time.Now()
		next.ServeHTTP(w, r)

		err := httph.ErrorGet(r)
		execTime := time.Since(start)

		if m.fromOptions.skipper(r) {
			return
		}

		// strings.Builder — эффективная сборка строки без промежуточных аллокаций.
		// Grow(48 + ...) резервирует память заранее: ~48 байт на метод + хвост сообщения.
		var mb strings.Builder
		mb.Grow(48 + len(r.RequestURI))
		mb.WriteString(r.Method)
		mb.WriteByte(' ')
		mb.WriteString(r.RequestURI)

		// Выбор уровня логирования
		var ev *zerolog.Event
		if err == nil {
			// Успешный запрос
			mb.WriteString(tailSuccess)
			ev = m.log.Debug()
		} else {
			// Ошибка
			mb.WriteString(tailFail)
			ev = m.log.Error()
		}

		ev.Err(err)
		ev.Ctx(r.Context())
		ev.Str("exec_time", execTime.String())
		ev.Str("client_ip", r.RemoteAddr)
		ev.Msg(mb.String())
	})
}

func NewMiddleware(opts ...Option) httph.Middleware {
	m := middleware{
		log: log.Logger,
	}
	m.fromOptions.skipper = defaultSkipper

	for _, opt := range opts {
		opt(&m)
	}

	return m.Callback
}

func defaultSkipper(_ *http.Request) bool {
	return false
}
