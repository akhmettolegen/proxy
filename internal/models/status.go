package models

type ErrorResponse struct {
	Error *ErrorStatus `json:"error"`
}

type ErrorStatus struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}
