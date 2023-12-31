package airazor

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/HuguesGuilleus/airazor/mock_transport"
	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	called := false

	config := Config{
		LimitBody: 5,
		Context:   context.Background(),
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

	request := &Request{
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
	}

	err := request.Fetch(&config)
	assert.NoError(t, err)
	assert.True(t, called)

	response := request.Response
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, mock_transport.HeaderServer, response.Header.Get("Server"))
	assert.Equal(t, mock_transport.Body[:5], string(response.Body))
	assert.Len(t, response.TestFails, 1)

	assert.True(t, called)
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
