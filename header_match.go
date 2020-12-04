// Package checkheadersplugin plugin.
package checkheadersplugin

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

//SingleHeader contains a single header keypair
type SingleHeader struct {
	Name     string `json:"name,omitempty"`
	Value    string `json:"value,omitempty"`
	Required *bool  `json:"required,omitempty"`
}

// Config the plugin configuration.
type Config struct {
	Headers []SingleHeader
}

// HeaderMatch demonstrates a HeaderMatch plugin.
type HeaderMatch struct {
	next    http.Handler
	headers []SingleHeader
	name    string
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Headers: []SingleHeader{},
	}
}

// New created a new HeaderMatch plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Headers) == 0 {
		return nil, fmt.Errorf("configuration incorrect, missing headers")
	}

	for _, vHeader := range config.Headers {
		if strings.TrimSpace(vHeader.Name) == "" {
			return nil, fmt.Errorf("configuration incorrect, missing header name")
		}
		if strings.TrimSpace(vHeader.Value) == "" {
			return nil, fmt.Errorf("configuration incorrect, missing header value")
		}
	}

	return &HeaderMatch{
		headers: config.Headers,
		next:    next,
		name:    name,
	}, nil
}

func (a *HeaderMatch) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, vHeader := range a.headers {
		reqHeaderVal := req.Header.Get(vHeader.Name)

		if vHeader.IsRequired() && reqHeaderVal != vHeader.Value {
			http.Error(rw, "Not allowed", http.StatusForbidden)
			return
		}

		if !vHeader.IsRequired() && reqHeaderVal != "" && reqHeaderVal != vHeader.Value {
			http.Error(rw, "Not allowed", http.StatusForbidden)
			return
		}
	}

	a.next.ServeHTTP(rw, req)
}

//IsRequired checks whether a header is mandatory in the request, defaults to 'true'
func (s *SingleHeader) IsRequired() bool {
	if s.Required == nil || *s.Required != false {
		return true
	}

	return false
}
