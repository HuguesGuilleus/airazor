package mock_transport

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMock(t *testing.T) {
	response, err := Transport.RoundTrip(httptest.NewRequest("GET", UrlOk, nil))
	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.Equal(t, Body, string(body))
}
