package server

import (
	"context"
	"encoding/json"
	"net/http"

	botUcase "github.com/drhidians/testbot/bot/usecase"
	"github.com/drhidians/testbot/server/ipchecker"
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

		r.Get("/", h.Start)
	})
	return r
}

func (h *botHandler) Start(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	u, err := h.s.Get(ctx)

	if err != nil {
		encodeError(ctx, err, w)
		return
	}

	var response = u

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Log("error", err)
		encodeError(ctx, err, w)
		return
	}
}
