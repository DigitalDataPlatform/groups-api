package router

import (
	"net/http"

	"github.com/go-chi/chi"
)

func NewConfigRouter(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Get Config"))
	})

	r.Get("/reload", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Reload Config"))
	})
}
