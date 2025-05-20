package example

// ExampleService is a service for interacting with the example API
//
//go:generate go run github.com/stijmetkii/validation-gen/cmd/gen
type ExampleService struct {
	Client *Client
	Cfg    *Config
}

// AnotherService is another example service
//
//go:generate go run github.com/strijmetkii/validation-gen/cmd/gen
type AnotherService struct {
	Logger  Logger
	Options Options
	Timeout int
}

// Logger is a simple logging interface
type Logger interface {
	Debug(msg string)
	Info(msg string)
	Error(msg string)
}

// Options represents configuration options
type Options struct {
	Region   string
	Endpoint string
}

// Client represents an API client
type Client struct {
	// Client fields
}

// Config represents client configuration
type Config struct {
	// Config fields
}
