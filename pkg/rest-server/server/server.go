package server

import (
	"github.com/gorilla/mux"
	"keeper/pkg/rest-server/handler"
	"net/http"
	"time"
)

func NewServer(h *handler.Handler, version string, port string) *http.Server {
	r := mux.NewRouter()

	// routes
	r.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		info := []byte(`{"version": "0.1.0"}`)

		_, err := w.Write(info)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
	}).Methods("GET")

	r.HandleFunc("/encrypt", h.Encrypt).Methods("POST")

	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return srv
}
