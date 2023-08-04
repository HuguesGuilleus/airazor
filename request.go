package airazor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/HuguesGuilleus/airazor/check"
)

type Request struct {
	parent *Collection

	Name string `json:"name,omitempty"`

	*Authorization `json:"authorization,omitempty"`

	Method string      `json:"method,omitempty"`
	URL    string      `json:"url,omitempty"`
	Header http.Header `json:"header,omitempty"`
	Body   string      `json:"body,omitempty"`

	Test string `json:"test,omitempty"`

	Response *Response `json:"response,omitempty"`
}

type Response struct {
	Name       string      `json:"name"`
	StatusCode int         `json:"StatusCode"`
	Header     http.Header `json:"Header"`
	Body       []byte      `json:"Body"`
	TestFails  []string    `json:"TestFails"`
}

type Config struct {
	NewContext func() context.Context
	http.RoundTripper
	LimitBody int64
}

func (config *Config) Fetch(request *Request) (*Response, error) {
	httpRequest, err := http.NewRequestWithContext(
		config.NewContext(),
		request.Method,
		request.URL,
		bytes.NewReader([]byte(request.Body)),
	)
	if err != nil {
		return nil, fmt.Errorf("Make request %q: %w", request.URL, err)
	}

	for key, values := range request.Header {
		for _, v := range values {
			httpRequest.Header.Add(key, v)
		}
	}

	if auth := request.getAuth(); auth != "" {
		httpRequest.Header.Set("Authorization", auth)
	}

	httpResponse, err := config.RoundTrip(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("Fetch %q: %w", request.URL, err)
	}
	defer httpResponse.Body.Close()

	body, err := io.ReadAll(io.LimitReader(httpResponse.Body, config.LimitBody))
	if err != nil {
		return nil, fmt.Errorf("Get all reponse body of %q: %w", request.URL, err)
	}

	response := &Response{
		StatusCode: httpResponse.StatusCode,
		Header:     httpResponse.Header,
		Body:       body,
	}

	response.test(request.Test)

	return response, nil
}

func (r *Response) test(src string) {
	if src == "" {
		return
	}
	r.TestFails = check.Run(src, map[string]any{
		"code": r.StatusCode,
		"text": func() string { return string(r.Body) },
	})
}

func (r *Request) getAuth() (auth string) {
	now := time.Now()

	auth = r.Authorization.Header(now)
	parent := r.parent
	for auth == "" && parent != nil {
		if parent.Authorization != nil && parent.Authorization.None {
			return
		}
		auth = parent.Authorization.Header(now)
		parent = parent.parent
	}
	return
}
