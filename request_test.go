package airazor

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.capgemini.com/hugues-guilleus/airazor/mock_transport"
	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	called := false

	config := Config{
		LimitBody:  5,
		NewContext: context.Background,
		RoundTripper: mock_transport.Inspector(func(request *http.Request) {
			called = true
			assert.Equal(t, "METHOD", request.Method)
			assert.Equal(t, []string{"V1", "V2"}, request.Header.Values("H"))
			assert.Equal(t, []string{"sésame"}, request.Header.Values("Authorization"))

			body, err := io.ReadAll(request.Body)
			assert.NoError(t, err)
			assert.Equal(t, "The body", string(body))
		}),
	}

	response, err := config.Fetch(&Request{
		parent: &Collection{
			Authorization: &Authorization{Raw: "sésame"},
		},
		Method: "METHOD",
		URL:    mock_transport.UrlOk,
		Header: http.Header{
			"H":             []string{"V1", "V2"},
			"Authorization": []string{"removed value"},
		},
		Body: "The body",
		Test: `assert(42, code);`,
	})
	assert.NoError(t, err)
	assert.True(t, called)

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, mock_transport.HeaderServer, response.Header.Get("Server"))
	assert.Equal(t, mock_transport.Body[:5], string(response.Body))
	assert.Len(t, response.TestFails, 1)

	assert.True(t, called)
}

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
