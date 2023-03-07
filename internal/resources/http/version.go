package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
)

const APIVersion = "v1"

type VersionResponse struct {
	API     string `json:"api"`
	Version string `json:"version"`
}

type VersionResource struct {
	Version string
}

func (vr VersionResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", vr.Get)

	return r
}

func (vr VersionResource) Get(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, VersionResponse{
		API:     APIVersion,
		Version: vr.Version,
	})
}
