package server

import (
	"context"
	"encoding/json"
	"net/http"

	tg "github.com/drhidians/testbot"
	botUcase "github.com/drhidians/testbot/bot/usecase"
	"github.com/drhidians/testbot/middleware/ipchecker"
	"github.com/go-chi/chi"
	kitlog "github.com/go-kit/kit/log"
)

type botHandler struct {
	s botUcase.Service

	logger kitlog.Logger
}

func (h *botHandler) router() chi.Router {
	r := chi.NewRouter()

	// Protected routes
	r.Group(func(r chi.Router) {

		r.Use(ipchecker.Check)

		r.Post("/", h.Update)
	})
	return r
}

func (h *botHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var upd tg.Update
	err := json.NewDecoder(r.Body).Decode(&upd)

	if err != nil {
		encodeError(ctx, err, w)
		return
	}

	err = h.s.Update(ctx, upd)

	if err != nil {
		encodeError(ctx, err, w)
		return
	}
}
