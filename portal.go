package main

import (
	"net/http"

	"github.com/apex/log"
	"github.com/go-chi/chi"
	"github.com/gobuffalo/envy"
)

func main() {
	port := ":" + envy.Get("PORT", "3000")
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World !"))
	})

	log.Infof("======== App running in %v mode ========", "stage")
	http.ListenAndServe(port, r)

}
