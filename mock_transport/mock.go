package mock_transport

import (
	"net/http"
	"net/http/httptest"
)

const Transport = transport(0)
const (
	Body = "Hello golang!"

	UrlOk   = "https://example.com/ok"
	UrlFail = "https://example.com/fail"
)

type transport uintptr

func (transport) RoundTrip(request *http.Request) (*http.Response, error) {
	response := httptest.NewRecorder()
	switch request.URL.String() {
	case UrlOk:
		response.Code = http.StatusOK
	case UrlFail:
		response.Code = http.StatusBadRequest
	default:
		response.Code = http.StatusNotFound
	}
	response.HeaderMap.Set("Server", "mockTransport")
	response.Write([]byte(Body))
	return response.Result(), nil
}
