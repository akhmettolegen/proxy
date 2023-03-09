package v1

import (
	"encoding/json"
	"github.com/akhmettolegen/proxy/internal/managers"
	"github.com/akhmettolegen/proxy/internal/managers/auth"
	"github.com/akhmettolegen/proxy/internal/managers/task"
	"github.com/akhmettolegen/proxy/internal/models"
	"github.com/akhmettolegen/proxy/internal/models/httperrors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
)

type TaskResource struct {
	TaskManager managers.TaskManager
	AuthManager *auth.AuthManager
}

func (rs TaskResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(auth.NewUserAccessCtx(rs.AuthManager.JWTKey()).ChiMiddleware)

	r.Group(func(r chi.Router) {

		r.Post("/", rs.taskCreate)
		r.Get("/{id}", rs.taskById)
	})

	return r
}

// @Tags taskCreate
// @Description Create task
// @Accept  json
// @Produce  json
// @Param body body models.TaskRequest true "Request"
// @Success 200 {object} models.TaskResponse
// @Failure 400 {object} httperrors.Response
// @Failure 401 {object} httperrors.Response
// @Failure 500 {object} httperrors.Response
// @Router /task [post]
func (rs TaskResource) taskCreate(w http.ResponseWriter, r *http.Request) {

	var req *models.TaskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}

	result, err := rs.TaskManager.TaskCreate(req)
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
func (rs TaskResource) taskById(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	if id == "" {
		_ = render.Render(w, r, httperrors.BadRequest(task.ErrInValidTaskId))
		return
	}

	result, err := rs.TaskManager.TaskById(id)
	if err != nil {
		_ = render.Render(w, r, httperrors.Internal(err))
		return
	}

	render.JSON(w, r, result)
}
