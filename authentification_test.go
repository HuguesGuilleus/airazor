package airazor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuthentificationEmpty(t *testing.T) {
	checkAuth(t, "", nil)
}
func TestAuthentificationNone(t *testing.T) {
	checkAuth(t, "", &Authorization{None: true})
}
func TestAuthentificationBasic(t *testing.T) {
	checkAuth(t, "Basic YWxhZGRpbjpzZXNhbWVPdXZyZVRvaQ==", &Authorization{
		Basic: &AuthorizationBasic{"aladdin", "sesameOuvreToi"},
	})
}
func TestAuthentificationBearer(t *testing.T) {
	checkAuth(t, "Bearer yolo", &Authorization{Bearer: "yolo"})
}
func TestAuthentificationRaw(t *testing.T) {
	checkAuth(t, "yolo", &Authorization{Raw: "yolo"})
}
func TestAuthentificationJWTNone(t *testing.T) {
	checkAuth(t, "Bearer eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJleHAiOjE2ODgxNjk2NjAsImlhdCI6MTY4ODE2OTYwMCwibmFtZSI6IkxhcmEgQ1JPRlQiLCJzdWIiOiJsYXJhQGNyb2Z0LmNvLnVrIn0.", &Authorization{JWT: &AuthorizationJWT{
		Body: map[string]any{
			"sub":  "lara@croft.co.uk",
			"name": "Lara CROFT",
		},
	}})
}
func TestAuthentificationJWTHS256(t *testing.T) {
	checkAuth(t, "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODgxNjk2NjAsImlhdCI6MTY4ODE2OTYwMCwibmFtZSI6IkxhcmEgQ1JPRlQiLCJzdWIiOiJsYXJhQGNyb2Z0LmNvLnVrIn0.QUvM2IWxphicBD-Vw0dkIKxn4qQ1sUVYbwdsgrrnvEY", &Authorization{JWT: &AuthorizationJWT{
		Body: map[string]any{
			"sub":  "lara@croft.co.uk",
			"name": "Lara CROFT",
		},
		Algo: "HS256",
		Key:  []byte("code"),
	}})
}
func TestAuthentificationJWTHS384(t *testing.T) {
	checkAuth(t, "Bearer eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODgxNjk2NjAsImlhdCI6MTY4ODE2OTYwMCwibmFtZSI6IkxhcmEgQ1JPRlQiLCJzdWIiOiJsYXJhQGNyb2Z0LmNvLnVrIn0.23aoUzlDPk7nuVaAAeI6kTRXxBW9NqJ8qu_EMltx1Tjkvf9yp0Q3FzMtPgPjFtbO", &Authorization{JWT: &AuthorizationJWT{
		Body: map[string]any{
			"sub":  "lara@croft.co.uk",
			"name": "Lara CROFT",
		},
		Algo: "HS384",
		Key:  []byte("code"),
	}})
}
func TestAuthentificationJWTHS512(t *testing.T) {
	checkAuth(t, "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODgxNjk2NjAsImlhdCI6MTY4ODE2OTYwMCwibmFtZSI6IkxhcmEgQ1JPRlQiLCJzdWIiOiJsYXJhQGNyb2Z0LmNvLnVrIn0.O4BUqnCTA00Fu5HQ8dJzL-23PpsxHVe1wOc34RjDjQ7Nh953Oq5iVqHPkYv3zBt6WjUorCXiuUD9MWawJkxkQA", &Authorization{JWT: &AuthorizationJWT{
		Body: map[string]any{
			"sub":  "lara@croft.co.uk",
			"name": "Lara CROFT",
		},
		Algo: "HS512",
		Key:  []byte("code"),
	}})
}
func checkAuth(t *testing.T, header string, auth *Authorization) {
	assert.Equal(t, header, auth.Header(time.Unix(1688169600, 0)))
}
