package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/YLysov0017/go_pet_n1/internal/config/storage"
	"github.com/YLysov0017/go_pet_n1/internal/lib/logger/api/response"
	"github.com/YLysov0017/go_pet_n1/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2@v2.49.2 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		resURL, err := urlGetter.GetURL(alias)
		switch {
		case errors.Is(err, storage.ErrURLNotFound):
			log.Info("url not found", slog.String("alias", alias))

			render.JSON(w, r, response.Error("not found"))

			return
		case err != nil:
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		log.Info("got url", slog.String("url", resURL))

		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
