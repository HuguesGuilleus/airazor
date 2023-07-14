package airazor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthentificationNone(t *testing.T) {
	assert.Equal(t, "", (&Authorization{None: true}).Header("foo"))
}
func TestAuthentificationBasic(t *testing.T) {
	assert.Equal(t, "Basic YWxhZGRpbjpzZXNhbWVPdXZyZVRvaQ==",
		(&Authorization{
			Basic: &AuthorizationBasic{"aladdin", "sesameOuvreToi"},
		}).Header("foo"))
}
func TestAuthentificationBearer(t *testing.T) {
	assert.Equal(t, "Bearer yolo",
		(&Authorization{Bearer: "yolo"}).Header("foo"))
}
func TestAuthentificationRaw(t *testing.T) {
	assert.Equal(t, "yolo",
		(&Authorization{Raw: "yolo"}).Header("foo"))
}
