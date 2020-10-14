package traefik_replace_response_code

import (
	"context"
	"fmt"
	"net/http"
)

type responseWriterWithStatusCode struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseWriterWithStatusCode) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// Config the plugin configuration.
type Config struct {
	inputCode  int `json:"inputCode,omitempty"`
	outputCode int `json:"outputCode,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		inputCode:  429,
		outputCode: 200,
	}
}

type Limiter struct {
	next       http.Handler
	inputCode  int
	outputCode int
	name       string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &Limiter{
		inputCode:  config.inputCode,
		outputCode: config.outputCode,
		next:       next,
		name:       name,
	}, nil
}

func (a *Limiter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	responseWriter := responseWriterWithStatusCode{rw, 200}
	fmt.Print("In Serve HTTP, calling next serve")
	a.next.ServeHTTP(&responseWriter, req)

	fmt.Sprintf("Status Code %d", responseWriter.statusCode)

	if responseWriter.statusCode == a.inputCode {
		responseWriter.WriteHeader(a.outputCode)
	}
}
