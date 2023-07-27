package airazor

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
)

type Collection struct {
	parent *Collection

	Name string

	Environments map[string]string
	Authorization

	Requests []Request
	Children []Collection
}

type Request struct {
	parent *Collection

	Name   string      `json:"name,omitempty"`
	Method string      `json:"method,omitempty"`
	URL    *url.URL    `json:"url,omitempty"`
	Header http.Header `json:"header,omitempty"`
	Body   string      `json:"body,omitempty"`

	Test string `json:"test,omitempty"`
}

func (r *Request) ID() string {
	data, _ := json.Marshal(r)
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}
