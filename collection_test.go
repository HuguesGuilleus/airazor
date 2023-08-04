package airazor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func creatFakeCollection() (*Collection, *Collection, *Collection) {
	root := &Collection{
		Requests: []*Request{{}, {}},
		Children: []*Collection{
			{
				Requests: []*Request{{
					Response: &Response{StatusCode: 400},
				}, {
					Response: &Response{StatusCode: 200},
				}},
				Children: []*Collection{
					{
						Requests: []*Request{{
							Response: &Response{StatusCode: 300},
						}, {
							Response: &Response{StatusCode: 500},
						}},
					},
				},
			},
		},
	}
	root.buildTree()
	child := root.Children[0]
	subchild := child.Children[0]
	return root, child, subchild
}

func TestCollectionBuildTree(t *testing.T) {
	root, child, subchild := creatFakeCollection()
	assert.Nil(t, root.parent)
	assert.Same(t, root, root.Requests[0].parent)
	assert.Same(t, root, root.Requests[1].parent)
	assert.Same(t, child, child.Requests[0].parent)
	assert.Same(t, child, child.Requests[1].parent)
	assert.Same(t, subchild, subchild.Requests[0].parent)
	assert.Same(t, subchild, subchild.Requests[1].parent)
}

func TestCollectionRemoveResponse(t *testing.T) {
	root, child, subchild := creatFakeCollection()
	root.removeResponse()
	assert.Nil(t, root.Requests[0].Response)
	assert.Nil(t, root.Requests[1].Response)
	assert.Nil(t, child.Requests[0].Response)
	assert.Nil(t, child.Requests[1].Response)
	assert.Nil(t, subchild.Requests[0].Response)
	assert.Nil(t, subchild.Requests[1].Response)
}
