package mwlogger

import (
	"net/http"
	"time"

	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With( // Будет выполнено до обработки запроса
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor) // Получить сведения об ответе

			t1 := time.Now()
			// Эта функция будет выполнена после всех остальных middleware-хэндлеров
			// и обработки запроса
			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			next.ServeHTTP(ww, r) // Передать управление следующему хэндлеру в цепочке
		}

		return http.HandlerFunc(fn)
	}
}