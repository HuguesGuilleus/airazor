package airazor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
	r := Request{Name: "yolo"}
	assert.Equal(t, "c01f6779d46365fa43a878480aabf0eda45adc0cb26516d2bc7befa92f09bffe", r.ID())
}
