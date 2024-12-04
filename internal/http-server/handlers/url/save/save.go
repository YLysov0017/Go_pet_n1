package save

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/YLysov0017/go_pet_n1/internal/config/storage"
	"github.com/YLysov0017/go_pet_n1/internal/lib/logger/api/response"
	"github.com/YLysov0017/go_pet_n1/internal/lib/logger/sl"
	"github.com/YLysov0017/go_pet_n1/internal/lib/random"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.49.2 --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave, alias string) (int64, error)
} // Переданный storage с этим методом автоматически удовлетворяет интерфейсу

func New(log *slog.Logger, urlSaver URLSaver, aliasLength int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRanomdString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		switch {
		case errors.Is(err, storage.ErrURLExists):
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, response.Error("url already exists"))

			return
		case errors.Is(err, storage.ErrAliasExists):
			log.Info("alias already exists", slog.String("url", req.URL))

			render.JSON(w, r, response.Error("alias already exists, please try again"))

			return
		case err != nil:
			log.Error("failed to add url", sl.Err(err)) // sl.Err - небезопасен для клиентской ошибки, возвращает тип хранилища

			render.JSON(w, r, response.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}
}
