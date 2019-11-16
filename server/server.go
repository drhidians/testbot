package server

import (
	"context"
	"encoding/json"
	"net/http"

	botUcase "github.com/drhidians/testbot/bot/usecase"
	userUcase "github.com/drhidians/testbot/user/usecase"
	"github.com/go-chi/chi"
	kitlog "github.com/go-kit/kit/log"
)

// Server holds the dependencies for a HTTP server.
type Server struct {
	Bot  botUcase.Service
	User userUcase.Service

	Logger kitlog.Logger

	router chi.Router
}

// New returns a new HTTP server.
func New(bot botUcase.Service, user userUcase.Service, logger kitlog.Logger, jwtToken string) *Server {
	s := &Server{
		Bot:    bot,
		User:   user,
		Logger: logger,
	}

	r := chi.NewRouter()

	r.Use(accessControl)

	r.Route("/bot/webhook", func(r chi.Router) {
		h := botHandler{s.Bot, s.Logger}
		r.Mount("/", h.router())
	})
	r.Route("/api/v1", func(r chi.Router) {
		h := apiHandler{s.User, s.Logger, jwtToken}
		r.Mount("/", h.router())
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
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
	case userUcase.ErrUserNotFound:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
