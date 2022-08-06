package manager

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gogoapps/providers"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

const (
	PORT=":8080"

	EnvPort = "PORT"
	EnvConcurrentRequests = "CONCURRENT_REQUESTS"
	EnvApiKey = "API_KEY"
	EnvProvider = "PROVIDER"
)

type Manager struct {
	Provider providers.Provider
	R chi.Router
	L *logrus.Logger
}

func NewManager() (*Manager, error) {
	m := &Manager{}

	for _, f := range []func() error{
		m.initLogger,
		m.initProvider,
		m.initRouter,
	} {
		err := f()
		if err != nil {
			logrus.Errorf("init manager methods failed: %v", err)
			return nil, fmt.Errorf("init manager methods failed: %v", err)
		}
	}

	return m, nil
}

func (m *Manager) initProvider() error {
	switch os.Getenv(EnvProvider) {
	case providers.NASA:
		concurrentRequests, err := strconv.Atoi(getEnvOrDefault(EnvConcurrentRequests))
		if err != nil {
			return err
		}

		m.Provider = &providers.Nasa{
			ApiKey: getEnvOrDefault(EnvApiKey),
			L: m.L,
			ConcurrentRequests: concurrentRequests,
		}
		return nil
	default:
		m.L.Errorf("could not get provider from env")
		return fmt.Errorf("could not get provider from env")
	}
}

func (m *Manager) initRouter() error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	m.R = r

	return nil
}

func (m *Manager) initLogger() error {
	m.L = logrus.New()

	return nil
}

func getEnvOrDefault(key string) (value string) {
	switch key {
	case EnvPort:
		value = PORT
		if os.Getenv(EnvPort) != "" {
			value = ":"+os.Getenv(EnvPort)
		}
	case EnvApiKey:
		value = providers.ApiKeyDefault
		if os.Getenv(EnvApiKey) != "" {
			value = os.Getenv(EnvApiKey)
		}
	case EnvConcurrentRequests:
		value = strconv.Itoa(providers.ConncurentRequests)
		if os.Getenv(EnvConcurrentRequests) != "" {
			value = os.Getenv(EnvConcurrentRequests)
		}
	case EnvProvider:
		value = providers.NASA
		if os.Getenv(EnvProvider) != "" {
			value = os.Getenv(EnvProvider)
		}
	}
	return
}
