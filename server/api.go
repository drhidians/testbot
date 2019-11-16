package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	kitlog "github.com/go-kit/kit/log"

)

type botHandler struct {
	s booking.Service

	logger kitlog.Logger
}

func (h *bookingHandler) router() chi.Router {
	r := chi.NewRouter()

	r.Route("/cargos", func(r chi.Router) {
		r.Post("/", h.bookCargo)
		r.Get("/", h.listCargos)
		r.Route("/{trackingID}", func(r chi.Router) {
			r.Get("/", h.loadCargo)
			r.Get("/request_routes", h.requestRoutes)
			r.Post("/assign_to_route", h.assignToRoute)
			r.Post("/change_destination", h.changeDestination)
		})

	})
	r.Get("/locations", h.listLocations)

	r.Method("GET", "/docs", http.StripPrefix("/booking/v1/docs", http.FileServer(http.Dir("booking/docs"))))

	return r
}

}
