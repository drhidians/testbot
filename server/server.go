package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/drhidians/testbot/bot"
	"github.com/drhidians/testbot/user"
	kitlog "github.com/go-kit/kit/log"
)

// Server holds the dependencies for a HTTP server.
type Server struct {
	Bot  bot.Service
	BotApi botapi.Service

	Logger kitlog.Logger

	router chi.Router
}

// New returns a new HTTP server.
func New(bot bot.Service,api api.Service, logger kitlog.Logger) *Server {
	s := &Server{
		Bot:  bot,
		Api: api
		Logger:   logger,
	}

	r := chi.NewRouter()

	r.Use(accessControl)

	r.Route("/bot/webhook", func(r chi.Router) {
		h := botHandler{s.Bot, s.Logger}
		r.Mount("/", h.router())
	})
	r.Route("/api", func(r chi.Router) {
		h := apiHandler{s.Api, s.Logger}
		r.Mount("/v1", h.router())
	})

	s.router = r

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case shipping.ErrUnknownCargo:
		w.WriteHeader(http.StatusNotFound)
	case tracking.ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
