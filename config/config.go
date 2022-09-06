package config

import (
	"github.com/spf13/viper"
)

const (
	PORT = "8080"
	ConcurrentRequests = 5
	Provider = "NASA"
	ApiKeyDefault = "DEMO_KEY"

	EnvPort               = "PORT"
	EnvConcurrentRequests = "CONCURRENT_REQUESTS"
	EnvApiKey             = "API_KEY"
	EnvProvider           = "PROVIDER"
)

type Config struct {
	Port string
	ConcurrentRequests int
	ApiKey string
	Provider string
}

func Init() (*Config, error) {
	viper.SetDefault(EnvPort, PORT)
	viper.SetDefault(EnvProvider, Provider)
	viper.SetDefault(EnvApiKey, ApiKeyDefault)
	viper.SetDefault(EnvConcurrentRequests, ConcurrentRequests)

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:               viper.GetString(EnvPort),
		ConcurrentRequests: viper.GetInt(EnvConcurrentRequests),
		ApiKey:             viper.GetString(EnvApiKey),
		Provider:           viper.GetString(EnvProvider),
	}, err
}

func (c Config) GetPort() string {
	return c.Port
}

func (c Config) GetConcurrentRequests() int {
	return c.ConcurrentRequests
}

func (c Config) GetApiKey() string {
	return c.ApiKey
}

func (c Config) GetProvider() string {
	return c.Provider
}