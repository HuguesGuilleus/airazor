package airazor

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
