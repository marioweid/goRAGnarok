package interfaces

type Server struct {
	BaseURL string
	APIKey  string
}

type Provider interface {
	Generate() string
}
