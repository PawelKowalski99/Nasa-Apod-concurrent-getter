package responses

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"-"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {

	render.Status(r, e.StatusCode)

	return nil
}

func NewErrResponse(errText string, statusCode int) render.Renderer {
	return &ErrResponse{
		Error:      errText,
		StatusCode: statusCode,
	}
}
