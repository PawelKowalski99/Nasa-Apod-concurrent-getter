package providers

import (
	"context"
	"net/http"
)

const ConncurentRequests = 5

type Provider interface {
	GetPictures(ctx context.Context) http.HandlerFunc
}