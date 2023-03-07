package clients

import "net/http"

type HttpClient interface {
	Request(method string, url string, headers map[string]string, body []byte) (*http.Response, error)
}
