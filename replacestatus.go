package traefik_replace_response_code

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
)


// Config the plugin configuration.
type Config struct {
	InputCode  int    `json:"inputCode,omitempty"`
	OutputCode int    `json:"outputCode,omitempty"`
	RemoveBody bool `json:"removeBody,omitempty"`

}

func CreateConfig() *Config {
	return &Config{
		InputCode:  429,
		OutputCode: 202,
		RemoveBody: false,
	}
}

type StatusCodeReplacer struct {
	next       http.Handler
	inputCode  int
	outputCode int
	removeBody bool
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	log.Printf("Configuring plugin replace-response-code with inputCode: %d, outputCode: %d, removeBody: %s", config.InputCode, config.OutputCode, config.RemoveBody)

	return &StatusCodeReplacer{
		inputCode:  config.InputCode,
		outputCode: config.OutputCode,
		removeBody: config.RemoveBody,
		next:       next,
	}, nil

}

func (a *StatusCodeReplacer) replacer() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		//recorder is used to delegate call. Its response will be used to create the correct ResponseWriter status code.
		recorder := httptest.NewRecorder()
		a.next.ServeHTTP(recorder, req)

		if recorder.Code == a.inputCode {
			rw.WriteHeader(a.outputCode)
			if !a.removeBody {
				_, _ = rw.Write(recorder.Body.Bytes())
			}
		}else{
			rw.WriteHeader(recorder.Code)
			_, _ = rw.Write(recorder.Body.Bytes())
		}

		for name, values := range recorder.Header(){
			rw.Header()[name] = values
		}



	})
}

func (a *StatusCodeReplacer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	a.replacer().ServeHTTP(rw,req)

}
