package http_client

import "net/http"

type HttpClient struct {
	httpClient *http.Client
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		httpClient: &http.Client{},
	}
}

func (hc *HttpClient) Do(req *http.Request) (*http.Response, error) {
	return hc.httpClient.Do(req)
}
