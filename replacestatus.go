package traefik_replace_response_code

import (
	"context"
	"log"
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
	InputCode  int `json:"inputCode,omitempty"`
	OutputCode int `json:"outputCode,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		InputCode:  429,
		OutputCode: 200,
	}
}

type Limiter struct {
	next       http.Handler
	inputCode  int
	outputCode int
	name       string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	log.Printf("Configuring plugin replace-response-code with inputCode: %d, outputCode: %d", config.InputCode, config.OutputCode)

	return &Limiter{
		inputCode:  config.InputCode,
		outputCode: config.OutputCode,
		next:       next,
		name:       name,
	}, nil
}

func (a *Limiter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	responseWriter := responseWriterWithStatusCode{rw, 200}
	log.Print("In Serve HTTP, calling next serve")
	a.next.ServeHTTP(&responseWriter, req)

	log.Printf("Status Code %t", responseWriter.statusCode == a.inputCode)

	if responseWriter.statusCode == a.inputCode {
		responseWriter.WriteHeader(a.outputCode)
		log.Printf("Status Code %d", responseWriter.statusCode)

	}
}
