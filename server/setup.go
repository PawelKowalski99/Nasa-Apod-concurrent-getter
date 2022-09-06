package server

import (
	"context"
	"github.com/PawelKowalski99/gogapps/config"
	"github.com/PawelKowalski99/gogapps/providers/nasa"
	"time"

	"github.com/PawelKowalski99/gogapps/providers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Config interface {
	GetPort() string
}

type Server struct {
	Config Config
	Provider providers.Provider
	R        chi.Router
}

func Init(c *config.Config) (*Server, error) {
	m := &Server{
		Config: c,
	}

	m.initProvider(c).initRouter().initRoutes()

	return m, nil
}

func (s *Server) initProvider(c nasa.Config) *Server {
	switch c.GetProvider() {
	default:
		s.Provider = &nasa.Nasa{
			Config: c,
		}
		return s
	}
}

func (s *Server) initRouter() *Server {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.ThrottleBacklog(1, 50, 60*time.Second))

	s.R = r

	return s
}

func (s *Server) initRoutes() *Server {
	s.R.Get("/pictures", s.Provider.GetPictures(context.Background()))
	return s
}