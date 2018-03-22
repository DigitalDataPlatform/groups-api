package router

import (
	"net/http"

	"github.com/apex/log"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	ddpportal "gitlab.adeo.com/ddp-portal-api"
	repository "gitlab.adeo.com/ddp-portal-api/api/repository"
)

type GroupRequest struct {
	*ddpportal.Group `json:"group"`
	User             string `json:"user,omitempty"`
}

func (g *GroupRequest) Bind(r *http.Request) error {
	return nil
}

type GroupResponse struct {
	*ddpportal.Group
}

func NewGroupResponse(group *ddpportal.Group) *GroupResponse {
	return &GroupResponse{Group: group}
}

func (rd *GroupResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type GroupListResponse []*GroupResponse

func NewGroupListResponse(groups []*ddpportal.Group) []render.Renderer {
	list := []render.Renderer{}
	for _, group := range groups {
		list = append(list, NewGroupResponse(group))
	}
	return list
}

type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

func NewGroupRouter(repo repository.GroupRepository) func(chi.Router) {
	return func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			data := &GroupRequest{}

			if err := render.Bind(r, data); err != nil {
				log.WithError(err).Error("Unable to parse input")
				render.Render(w, r, ErrInvalidRequest(err))
				return
			}

			group := ddpportal.NewGroup(data.Group.Name, data.User)
			err := repo.Insert(group)
			if err != nil {
				log.WithError(err).Error("Unable to insert group")
				render.Status(r, http.StatusInternalServerError)
				return
			}

			render.Status(r, http.StatusCreated)
			render.Render(w, r, NewGroupResponse(&group))
		})
	}
}
