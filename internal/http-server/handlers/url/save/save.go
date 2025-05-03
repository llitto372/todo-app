package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"todo/internal/config/storage"
	"todo/internal/lib/api/responce"
	"todo/internal/lib/logger/sl"
	"todo/internal/lib/random"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	responce.Response `json:"url" validate:"required, url"`
	Alias             string `json:"alias"`
}
type URLSaver interface {
	SaveURL(URLtoSave string, alias string) error
}

const aliasLength = 6

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(slog.String("op", op), slog.String("requested_id", middleware.GetReqID(r.Context())))
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to parse request", sl.Err(err))

			render.JSON(w, r, responce.Error("failed to decode request"))

			return
		}

		log.Info("request body decodecd", slog.Any("req", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("failed to validate request", sl.Err(err))

			render.JSON(w, r, responce.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		err = urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))
			render.JSON(w, r, responce.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to save url", sl.Err(err))

			render.JSON(w, r, responce.Error("failed to save url"))
			return
		}
		log.Info("url added", slog.String("url", req.URL))

		render.JSON(w, r, Response{
			Response: responce.OK(),
			Alias:    alias,
		})
	}

}
