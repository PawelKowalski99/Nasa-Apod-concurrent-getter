package providers

import (
	"context"
	"net/http"
)

type Provider interface {
	GetPictures(ctx context.Context) http.HandlerFunc
}