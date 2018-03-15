package main

import (
	"net/http"

	"github.com/apex/log"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gobuffalo/envy"
	ddpportal "gitlab.adeo.com/ddp-portal-api"
)

func main() {
	port := ":" + envy.Get("PORT", "3000")
	oauthURL := envy.Get("OAUTH_SERVER", "https://yradeooauth.corp.leroymerlin.com/oauth-server")
	oauthClientID := envy.Get("CLIENT_ID", "CORP-DDP-PORTAL-DEV")
	oauthClientSecret := envy.Get("CLIENT_SECRET", "imtj8d0572dtg7glwb93ss")

	r := chi.NewRouter()

	adeoOauthProvider := ddpportal.NewOauthAdeoProvider(oauthURL, oauthClientID, oauthClientSecret)
	oauth := ddpportal.NewOauthAuthorizerMiddleware(adeoOauthProvider)

	r.Use(middleware.RequestID)
	r.Use(oauth)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/config", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Get Config"))
		})

		r.Get("/reload", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Reload Config"))
		})
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World !"))
	})

	log.Infof("======== App running in %v mode ========", "stage")
	http.ListenAndServe(port, r)

}
