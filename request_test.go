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

func TestRequestGetAuth(t *testing.T) {
	root := &Collection{
		Authorization: &Authorization{Raw: "y1"},
		Children: []*Collection{{
			Authorization: &Authorization{None: true},
			Children: []*Collection{{
				Authorization: &Authorization{Raw: "y3"},
				Requests: []*Request{{
					Authorization: &Authorization{Raw: "y4"},
				}},
			}},
		}},
	}
	root.buildTree()

	assert.Equal(t, "y4", root.Children[0].Children[0].Requests[0].getAuth())
	root.Children[0].Children[0].Requests[0].Authorization = nil
	assert.Equal(t, "y3", root.Children[0].Children[0].Requests[0].getAuth())
	root.Children[0].Children[0].Authorization = nil
	assert.Equal(t, "", root.Children[0].Children[0].Requests[0].getAuth())
	root.Children[0].Authorization = nil
	assert.Equal(t, "y1", root.Children[0].Children[0].Requests[0].getAuth())
	root.Authorization = nil
	assert.Equal(t, "", root.Children[0].Children[0].Requests[0].getAuth())
}
