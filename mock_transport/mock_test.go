package mock_transport

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMock(t *testing.T) {
	response, err := Transport.RoundTrip(httptest.NewRequest("GET", UrlOk, nil))
	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, []string{HeaderServer}, response.Header.Values("Server"))

	body, err := io.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.Equal(t, Body, string(body))
}

func TestInspector(t *testing.T) {
	globalRequest := httptest.NewRequest("GET", UrlOk, nil)

	called := 0
	Inspector(func(argRequest *http.Request) {
		called++
		assert.Same(t, globalRequest, argRequest)
	}).RoundTrip(globalRequest)

	assert.Equal(t, called, 1)
}
