package router

import (
	"net/http"

	"github.com/apex/log"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	repository "gitlab.adeo.com/ddp-portal-api/api/repository"
)

func NewMeRouter(repo repository.GroupRepository) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("My profil"))
		})

		r.Get("/groups", func(w http.ResponseWriter, r *http.Request) {
			groups, err := repo.GetAll()
			if err != nil {
				log.WithError(err).Error("Unable to read groups")
				render.Status(r, http.StatusInternalServerError)
				return
			}
			render.RenderList(w, r, NewGroupListResponse(groups))
		})
	}

}
