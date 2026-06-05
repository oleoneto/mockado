// Package wrapping Go's native http client.
package httpclient

import (
	"net/http"
	"sync"
	"time"
)

var (
	client     *http.Client
	clientOnce sync.Once
)

var _ HTTPClient = &http.Client{}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Example usage:
//
//	 url, _ := url.Parse("http://google.com/search?q=golang")
//
//	 request := http.Request{
//		  Method: "GET"
//		  URL:    url,
//		  Header: map[string][]string{"User-Agent": "curl"},
//		  Body: &RequestBody{Data: []byte()},
//	 }
//
//	 response, err := HTTPClient.Do(&request)
func Client() *http.Client { return client }

func WithTimeout(t time.Duration) *http.Client {
	client.Timeout = t
	return client
}

func init() {
	clientOnce.Do(func() {
		client = &http.Client{}
	})
}
