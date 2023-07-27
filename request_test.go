package airazor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
	r := Request{Name: "yolo"}
	assert.Equal(t, "c01f6779d46365fa43a878480aabf0eda45adc0cb26516d2bc7befa92f09bffe", r.ID())
}

func TestCollectionBuildTree(t *testing.T) {
	root := &Collection{
		Requests: []*Request{{}, {}},
		Children: []*Collection{
			{
				Requests: []*Request{{}, {}},
				Children: []*Collection{
					{
						Requests: []*Request{{}, {}},
					},
				},
			},
		},
	}
	child := root.Children[0]
	subchild := child.Children[0]

	root.buildTree()

	assert.Nil(t, root.parent)
	assert.Same(t, root, root.Requests[0].parent)
	assert.Same(t, root, root.Requests[1].parent)
	assert.Same(t, child, child.Requests[0].parent)
	assert.Same(t, child, child.Requests[1].parent)
	assert.Same(t, subchild, subchild.Requests[0].parent)
	assert.Same(t, subchild, subchild.Requests[1].parent)
}
