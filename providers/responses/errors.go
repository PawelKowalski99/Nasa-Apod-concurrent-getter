package responses

import (
	"github.com/go-chi/render"
	"net/http"
)

type Error struct {
	HTTPStatusCode int    `json:"http_status_code"`
	StatusText     string `json:"status_text"`
	ErrorText      string `json:"error_text,omitempty"`
}

type ErrResponse struct {
	Error Error `json:"error"`
}
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.Error.HTTPStatusCode)

	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{Error: Error{
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     http.StatusText(http.StatusBadRequest),
		ErrorText:      err.Error(),
	}}
}

func ErrStatusNotOk(status int, errText string) render.Renderer {
	return &ErrResponse{Error: Error{
		HTTPStatusCode: status,
		StatusText:     http.StatusText(status),
		ErrorText:      errText,
	}}
}

func ErrParseBody(status int, err error) render.Renderer {
	return &ErrResponse{Error: Error{
		HTTPStatusCode: status,
		StatusText:     http.StatusText(status),
		ErrorText:      err.Error(),
	}}
}

func ErrConcurrentlyRun(status int, err error) render.Renderer {
	return &ErrResponse{Error: Error{
		HTTPStatusCode: status,
		StatusText:     http.StatusText(status),
		ErrorText:      err.Error(),
	}}
}