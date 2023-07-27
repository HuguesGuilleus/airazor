package airazor

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

type Collection struct {
	parent *Collection

	Name string

	*Authorization `json:"authorization,omitempty"`

	Requests []*Request
	Children []*Collection
}

type Request struct {
	parent *Collection

	Name string `json:"name,omitempty"`

	*Authorization `json:"authorization,omitempty"`

	Method string      `json:"method,omitempty"`
	URL    *url.URL    `json:"url,omitempty"`
	Header http.Header `json:"header,omitempty"`
	Body   string      `json:"body,omitempty"`

	Test string `json:"test,omitempty"`
}

// Get the ID of the request based on the
func (r *Request) ID() string {
	data, _ := json.Marshal(r)
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// Recursively recursively .parent of sub elements.
func (c *Collection) buildTree() {
	for _, child := range c.Children {
		child.parent = c
		child.buildTree()
	}
	for _, request := range c.Requests {
		request.parent = c
	}
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
