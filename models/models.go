package models

// GatewayRequest type....
type GatewayRequest struct {
	Path    string              `json:"path"`
	Params  map[string][]string `json:"multiValueParams"`
	Headers map[string]string   `json:"headers"`
	Body    string              `json:"body"`
}

// Response type ...
type Response struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}
