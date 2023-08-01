package mock_transport

import (
	"net/http"
	"net/http/httptest"
)

const Transport = transport(0)
const (
	HeaderServer = "mockTransport"

	Body = "Hello golang!"

	UrlOk   = "https://example.com/ok"
	UrlFail = "https://example.com/fail"
)

type transport uintptr

func (transport) RoundTrip(request *http.Request) (*http.Response, error) {
	response := httptest.NewRecorder()
	response.HeaderMap.Set("Server", HeaderServer)

	switch request.URL.String() {
	case UrlOk:
		response.Code = http.StatusOK
	case UrlFail:
		response.Code = http.StatusBadRequest
	default:
		response.Code = http.StatusNotFound
	}

	response.Write([]byte(Body))
	return response.Result(), nil
}

// Just check the request, then execute Transport.RoundTrip()
type Inspector func(*http.Request)

func (f Inspector) RoundTrip(request *http.Request) (*http.Response, error) {
	f(request)
	return Transport.RoundTrip(request)
}
