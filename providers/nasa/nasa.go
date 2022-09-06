package nasa



type Config interface {
	GetApiKey() string
	GetConcurrentRequests() int
	GetProvider() string
}

type Nasa struct {
	Config Config
}
