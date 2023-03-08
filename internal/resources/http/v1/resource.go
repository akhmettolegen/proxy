package v1

import (
	"encoding/json"
	"github.com/akhmettolegen/proxy/internal/managers"
	"github.com/akhmettolegen/proxy/internal/managers/proxy"
	"github.com/akhmettolegen/proxy/internal/models"
	"github.com/akhmettolegen/proxy/internal/models/httperrors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
)

type ProxyResource struct {
	ProxyManager managers.ProxyManager
}

func (rs ProxyResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {

		r.Post("/", rs.proxyRequest)
		r.Get("/{id}", rs.taskById)
	})

	return r
}

// @Tags proxyRequest
// @Description Proxy request
// @Accept  json
// @Produce  json
// @Param body body models.ProxyRequest true "Request"
// @Success 200 {object} models.ProxyResponse
// @Failure 400 {object} httperrors.Response
// @Failure 401 {object} httperrors.Response
// @Failure 500 {object} httperrors.Response
// @Router /task [post]
func (rs ProxyResource) proxyRequest(w http.ResponseWriter, r *http.Request) {

	var req *models.ProxyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}

	result, err := rs.ProxyManager.ProxyRequest(req)
	if err != nil {
		_ = render.Render(w, r, httperrors.Internal(err))
		return
	}

	render.JSON(w, r, result)
}

// @Tags taskById
// @Description Get task by id
// @Accept  json
// @Produce  json
// @Param id path string true "Task id"
// @Success 200 {object} models.Task
// @Failure 400 {object} httperrors.Response
// @Failure 401 {object} httperrors.Response
// @Failure 500 {object} httperrors.Response
// @Router /task/{id} [get]
func (rs ProxyResource) taskById(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	if id == "" {
		_ = render.Render(w, r, httperrors.BadRequest(proxy.ErrInValidTaskId))
		return
	}

	result, err := rs.ProxyManager.TaskById(id)
	if err != nil {
		_ = render.Render(w, r, httperrors.Internal(err))
		return
	}

	render.JSON(w, r, result)
}
