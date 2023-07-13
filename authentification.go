package airazor

import (
	"encoding/base64"
	"net/url"
)

// One of the following fields
type Authorization struct {
	None   bool
	Basic  *url.Userinfo
	Bearer string
	Raw    string
	JWT    *AuthorizationJWT
}

type AuthorizationJWT struct {
	Jose map[string]any
	Body map[string]any

	// HS256
	// HS384
	// HS512
}

// Generate the string to be used in HTTP headers. If the authentification is not defined,
func (auth *Authorization) Header(previous string) string {
	if auth.None {
		return ""
	}
	if basic := auth.Basic; basic != nil {
		return "Basic " +
			base64.URLEncoding.EncodeToString([]byte(basic.String()))
	}
	if bearer := auth.Bearer; bearer != "" {
		return "Bearer " + bearer
	}
	if raw := auth.Raw; raw != "" {
		return raw
	}

	return previous
}
