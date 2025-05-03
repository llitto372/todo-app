package redirect

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"todo/internal/config/storage"
	"todo/internal/lib/api/responce"
	"todo/internal/lib/logger/sl"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, URLgetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

		}
		resURL, err := URLgetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found")

			render.JSON(w, r, responce.Error("not found"))

			return
		}
		if err != nil {
			log.Error("internal error", sl.Err(err))

			render.JSON(w, r, responce.Error("internal error"))

			return
		}
		log.Info("success", resURL)

		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
