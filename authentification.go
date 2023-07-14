package airazor

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"hash"
	"time"
)

// One of the following fields
type Authorization struct {
	None   bool
	Basic  *AuthorizationBasic
	Bearer string
	Raw    string
	JWT    *AuthorizationJWT
}

type AuthorizationBasic struct {
	User, Password string
}

type AuthorizationJWT struct {
	Jose map[string]any
	Body map[string]any

	// Possibles values: HS256|HS384|HS512
	// Else value will produce a none algo signature
	Algo string
	Key  []byte
}

// Generate the string to be used in HTTP headers. If the authentification is not defined,
func (auth *Authorization) Header(previous string, now time.Time) string {
	if auth.None {
		return ""
	}
	if basic := auth.Basic; basic != nil {
		code := []byte(basic.User + ":" + basic.Password)
		return "Basic " + base64.URLEncoding.EncodeToString(code)
	}
	if bearer := auth.Bearer; bearer != "" {
		return "Bearer " + bearer
	}
	if raw := auth.Raw; raw != "" {
		return raw
	}
	if jwt := auth.JWT; jwt != nil {
		algo := jwtCanonizeAlgo(jwt.Algo)
		joseB, _ := json.Marshal(mergeMapAny(
			map[string]any{
				"alg": algo,
				"typ": "JWT",
			},
			jwt.Jose,
		))
		jose := base64.RawURLEncoding.EncodeToString(joseB)

		bodyB, _ := json.Marshal(mergeMapAny(
			map[string]any{
				"iat": now.Unix(),
				"exp": now.Add(time.Minute).Unix(),
			},
			jwt.Body,
		))
		body := base64.RawURLEncoding.EncodeToString(bodyB)

		signature := ""
		switch algo {
		case "none":
		case "HS256":
			signature = hamacAndBase64(sha256.New, jwt.Key, jose, ".", body)
		case "HS384":
			signature = hamacAndBase64(sha512.New384, jwt.Key, jose, ".", body)
		case "HS512":
			signature = hamacAndBase64(sha512.New, jwt.Key, jose, ".", body)
		}

		return "Bearer " + jose + "." + body + "." + signature
	}

	return previous
}

// Merge the 2 map in a new one.
// The a values overwrite b values.
func mergeMapAny(a, b map[string]any) (r map[string]any) {
	r = make(map[string]any, len(a)+len(b))
	for k, v := range b {
		r[k] = v
	}
	for k, v := range a {
		r[k] = v
	}
	return
}

func jwtCanonizeAlgo(algo string) string {
	switch algo {
	case "HS256":
		return "HS256"
	case "HS384":
		return "HS384"
	case "HS512":
		return "HS512"
	default:
		return "none"
	}
}

func hamacAndBase64(h func() hash.Hash, key []byte, datas ...string) string {
	hasher := hmac.New(h, key)
	for _, d := range datas {
		hasher.Write([]byte(d))
	}
	return base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))
}
