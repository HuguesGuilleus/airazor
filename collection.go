package airazor

import (
	"encoding/json"
	"sync"
)

type Collection struct {
	parent *Collection

	Name string `json:"name"`

	*Authorization `json:"authorization,omitempty"`

	Requests []*Request    `json:"requests"`
	Children []*Collection `json:"children"`
}

func ParseCollection(data []byte) (*Collection, error) {
	c := &Collection{}
	if err := json.Unmarshal(data, c); err != nil {
		return nil, err
	}
	c.buildTree()
	c.RemoveResponse()
	return c, nil
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

// Remove all response object of all requests.
func (c *Collection) RemoveResponse() {
	for _, child := range c.Children {
		child.RemoveResponse()
	}
	for _, request := range c.Requests {
		request.Error = ""
		request.Response = nil
	}
}

func (c *Collection) Fetch(config *Config) {
	wg := sync.WaitGroup{}
	defer wg.Wait()
	c.fetch(config, &wg)
}

func (c *Collection) fetch(config *Config, wg *sync.WaitGroup) {
	for _, child := range c.Children {
		child.fetch(config, wg)
	}
	wg.Add(len(c.Requests))
	for _, request := range c.Requests {
		go func(r *Request) {
			defer wg.Done()
			r.Fetch(config)
		}(request)
	}
}
