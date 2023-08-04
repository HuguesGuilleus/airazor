package airazor

import "sync"

type Collection struct {
	parent *Collection

	Name string `json:"name"`

	*Authorization `json:"authorization,omitempty"`

	Requests []*Request    `json:"requests"`
	Children []*Collection `json:"children"`
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
func (c *Collection) removeResponse() {
	for _, child := range c.Children {
		child.removeResponse()
	}
	for _, request := range c.Requests {
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
