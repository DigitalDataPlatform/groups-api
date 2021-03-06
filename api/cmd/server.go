package main

import (
	"net/http"

	"github.com/dgraph-io/badger"

	"github.com/apex/log"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gobuffalo/envy"
	ddpportal "gitlab.adeo.com/ddp-portal-api"
	"gitlab.adeo.com/ddp-portal-api/api/repository"
	"gitlab.adeo.com/ddp-portal-api/api/router"
)

func InitBadger() (*badger.DB, error) {
	opts := badger.DefaultOptions
	opts.Dir = "/tmp/badger"
	opts.ValueDir = "/tmp/badger"
	return badger.Open(opts)
}

func main() {
	port := ":" + envy.Get("PORT", "3000")

	authMiddleware := ddpportal.NewJwtAuthMiddleware()

	db, err := InitBadger()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	groupRepo := repository.NewBadgerGroupRepository(db)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(authMiddleware)
	r.Use(middleware.Recoverer)

	r.Get("/portal", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World !"))
	})

	r.Route("/portal/config", router.NewConfigRouter)
	r.Route("/portal/groups", router.NewGroupRouter(groupRepo))
	r.Route("/portal/me", router.NewMeRouter(groupRepo))

	log.Infof("======== App running in %v mode ========", "stage")
	http.ListenAndServe(port, r)

}
