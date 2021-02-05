package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
)

const (
	// TG ...
	TG = "/tg"
	// DEV ...
	DEV = "/dev"
)

// Handler -  Точка входа
func Handler(ctx context.Context, request *GatewayRequest) (*Response, error) {
	if request == nil {
		return nil, fmt.Errorf("nil request")
	}

	ID := getParam(request, "id")

	switch request.Path {
	case TG:
		return &Response{
			StatusCode: 200,
			Body:       tgForm(ID),
		}, nil
	case DEV:
		return &Response{
			StatusCode: 200,
			Body:       "this DEV" + ID + " " + request.Path,
		}, nil
	default:
		return &Response{
			StatusCode: 200,
			Body:       "null",
		}, nil
	}
}

// ViewData ...
type ViewData struct {
	ID     string
	Wallet string
}

func tgForm(ID string) string {
	data := ViewData{
		ID:     "TG-" + ID,
		Wallet: os.Getenv("YOO_WALLET"),
	}
	t, err := template.ParseFiles("form-tg.html")
	if err != nil {
		return err.Error()
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return err.Error()
	}
	return tpl.String()
}

// GatewayRequest type....
type GatewayRequest struct {
	Path    string              `json:"path"`
	Params  map[string][]string `json:"multiValueParams"`
	Headers map[string]string   `json:"headers"`
}

// Response type ...
type Response struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}

func getParam(r *GatewayRequest, name string) string {
	values := r.Params[name]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func getMultiparam(r *GatewayRequest, name string) []string {
	var res []string
	for _, value := range r.Params[name] {
		if len(value) > 0 {
			res = append(res, value)
		}
	}
	return res
}
