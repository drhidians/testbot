package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/drhidians/testbot/server/jwtauth"
	userUcase "github.com/drhidians/testbot/user/usecase"
	"github.com/go-chi/chi"
	kitlog "github.com/go-kit/kit/log"
)

var tokenAuth *jwtauth.JWTAuth

type apiHandler struct {
	s userUcase.Service

	logger   kitlog.Logger
	jwtToken string
}

func (h *apiHandler) router() chi.Router {
	r := chi.NewRouter()

	r.Get("/bot", h.GetBot)

	tokenAuth = jwtauth.New("HS256", []byte(h.jwtToken), nil)

	// Protected routes
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens

		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/user", h.GetUser)
	})
	return r
}

func (h *apiHandler) GetBot(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	b, err := h.s.GetBot(ctx)

	if err != nil {
		encodeError(ctx, err, w)
		return
	}

	var response = b

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Log("error", err)
		encodeError(ctx, err, w)
		return
	}
}

func (h *apiHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	_, claims, _ := jwtauth.FromContext(r.Context())
	id := claims["id"].(int64)

	u, err := h.s.GetByID(ctx, id)

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
