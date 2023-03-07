package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"github.com/akhmettolegen/proxy/internal/clients"
	"log"
	"net/http"
)

type Client struct {
	ctx context.Context
}

func NewClient(ctx context.Context) clients.HttpClient {
	return &Client{
		ctx: ctx,
	}
}
func (c *Client) Request(method string, url string, headers map[string]string, body []byte) (*http.Response, error) {

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Println("[ERROR] HTTP Create NewRequest error:", err)
		return nil, err
	}

	req = req.WithContext(c.ctx)

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("[ERROR] HTTP request error:", err)
		return nil, err
	}

	return resp, err
}
