package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	cache "github.com/drhidians/testbot/middleware/http-cache"
	"github.com/drhidians/testbot/middleware/http-cache/memory"
	"github.com/drhidians/testbot/middleware/jwtauth"
	userUcase "github.com/drhidians/testbot/user/usecase"
	"github.com/go-chi/chi"
	kitlog "github.com/go-kit/kit/log"
)

type apiHandler struct {
	s userUcase.Service

	logger kitlog.Logger
}

func (h *apiHandler) router() chi.Router {
	r := chi.NewRouter()

	r.Route("/media/{fileID}", func(r chi.Router) {

		memcached, err := memory.NewAdapter(
			memory.AdapterWithAlgorithm(memory.LRU),
			memory.AdapterWithCapacity(10000),
		)

		if err != nil {
			panic(err)
		}

		cacheClient, err := cache.NewClient(
			cache.ClientWithAdapter(memcached),
			cache.ClientWithTTL(24*60*time.Hour),
			cache.ClientWithRefreshKey("opn"),
		)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		r.Use(cacheClient.Middleware)
		r.Get("/", h.GetFile)

	})

	r.Get("/bot", h.GetBot)

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
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

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
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	_, claims, _ := jwtauth.FromContext(r.Context())

	id := int(claims["id"].(float64))

	u, err := h.s.GetByTelegramID(ctx, id)

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

func (h *apiHandler) GetFile(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	fileID := chi.URLParam(r, "fileID")

	//TO DO need to implement error
	if fileID == "" {
		encodeError(ctx, errors.New("Bad Request"), w)
		return
	}

	fileB, err := h.s.GetAvatar(ctx, fileID)
	if err != nil {
		encodeError(ctx, err, w)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename=avatar.png")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	http.ServeContent(w, r, "avatar.png", time.Now(), bytes.NewReader(fileB))
}

func downloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
