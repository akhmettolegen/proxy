package models

type ProxyRequest struct {
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    interface{}       `json:"body"`
}

type ProxyResponse struct {
	Id string `json:"id"`
}
