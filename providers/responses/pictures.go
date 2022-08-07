package responses

import (
	"github.com/go-chi/render"
	"net/http"
)

type PictureUrlsResponse struct {
	Urls []string `json:"urls"`
}

func (p *PictureUrlsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusOK)

	return nil
}

func PictureUrls(urls []string) render.Renderer {
	return &PictureUrlsResponse{Urls: urls}
}

type Picture struct {
	Copyright   string `json:"copyright"`
	Date        string `json:"date"`
	Explanation string `json:"explanation"`
	Title       string `json:"title"`
	Url         string `json:"url"`
	Hdurl       string `json:"hdurl"`
	MediaType   string `json:"media_type"`
	ServiceVersion string `json:"service_version"`
}

func (p *Picture) Render(w http.ResponseWriter, r *http.Request) error {
	err := render.DecodeJSON(r.Body, &p)
	if err != nil {
		return err
	}

	return nil
}

