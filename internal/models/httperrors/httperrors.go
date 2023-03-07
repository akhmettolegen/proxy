package httperrors

import (
	"net/http"

	"github.com/go-chi/render"
)

type Response struct {
	Err            error    `json:"-"`
	HTTPStatusCode int      `json:"-"`
	ErrorMessage   *Details `json:"error"`
	Validation     []string `json:"validation,omitempty"`
}

type Details struct {
	StatusText  string `json:"status"`
	AppCode     int64  `json:"code,omitempty"`
	MessageText string `json:"message,omitempty"`
}

func (e *Response) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func Internal(err error) render.Renderer {
	return &Response{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		ErrorMessage: &Details{
			AppCode:     http.StatusInternalServerError,
			StatusText:  "Internal Server Error",
			MessageText: err.Error(),
		},
	}
}

func BadRequest(err error) render.Renderer {
	return &Response{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		ErrorMessage: &Details{
			AppCode:     http.StatusBadRequest,
			StatusText:  http.StatusText(http.StatusBadRequest),
			MessageText: err.Error(),
		},
	}
}
