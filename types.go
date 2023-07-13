package airazor

import (
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

	Name   string
	Method string
	URL    *url.URL
	Header http.Header
	Body   string

	Test string
}
