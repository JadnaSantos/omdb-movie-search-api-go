package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github/JadnaSantos/omdb-movie-search-api-go.git/omdb"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewHandler (apikey string) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer) 
	r.Use(middleware.RequestID) 
	r.Use(middleware.Logger) 
	r.Use(jsonMiddleware)

	r.Get("/", handleSearchMovie(apikey))

	return r
}

type Respose struct {
	Error string `json:"error,omitempty"`
	Data any `json:"data,omitempty"`
}


func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func sendJSON (w http.ResponseWriter, resp Respose, status int) {
	data, err := json.Marshal(resp) 
	if err != nil {
		slog.Error("failed to marshal json data", "erro", err)
		sendJSON(
			w, 
			Respose{Error: "something went wrong"},
			http.StatusInternalServerError,
		)

		return 
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("failed to write response to client", "error", err)
		return 
	}
}

func handleSearchMovie (apiKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("s")
		res, err := omdb.Search(apiKey, search)

		if err != nil {
			sendJSON(
				w,
				Respose{Error: "something wrong with omdb"},
				http.StatusBadGateway,
			)
			return
		}

		sendJSON(w, Respose{Data: res}, http.StatusOK)
	}
}