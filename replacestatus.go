package traefik_replace_response_code

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
)


// Config the plugin configuration.
type Config struct {
	InputCode  int `json:"inputCode,omitempty"`
	OutputCode int `json:"outputCode,omitempty"`
	OutputBody *string `json:"outputBody,omitempty"`

}

func CreateConfig() *Config {
	return &Config{
		InputCode:  429,
		OutputCode: 202,
		OutputBody: nil,
	}
}

type StatusCodeReplacer struct {
	next       http.Handler
	inputCode  int
	outputCode int
	outputBody *string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	log.Printf("Configuring plugin replace-response-code with inputCode: %d, outputCode: %d, outputBody: %s", config.InputCode, config.OutputCode, config.OutputBody)

	if config.OutputBody != nil {
		return &StatusCodeReplacer{
			inputCode:  config.InputCode,
			outputCode: config.OutputCode,
			outputBody: config.OutputBody,
			next:       next,
		}, nil
	}else{
		return &StatusCodeReplacer{
			inputCode:  config.InputCode,
			outputCode: config.OutputCode,
			next:       next,
		}, nil
	}

}

func (a *StatusCodeReplacer) replacer() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		//recorder is used to delegate call. Its response will be used to create the correct ResponseWriter status code.
		recorder := httptest.NewRecorder()
		a.next.ServeHTTP(recorder, req)

		replaceBody := false

		if recorder.Code == a.inputCode {
			rw.WriteHeader(a.outputCode)
			if a.outputBody != nil {
				replaceBody = true
			}
			replaceBody = true
		}else{
			rw.WriteHeader(recorder.Code)

		}

		for name, values := range recorder.Header(){
			rw.Header()[name] = values
		}

		if replaceBody {
			_, _ = rw.Write([]byte(*a.outputBody))
		}else{
			_, _ = rw.Write(recorder.Body.Bytes())
		}

	})
}

func (a *StatusCodeReplacer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	a.replacer().ServeHTTP(rw,req)

}
