package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/morzik45/dnb-payment-form/models"
	"github.com/morzik45/dnb-payment-form/notifications"
)

const (
	// TG ...
	TG = "/tg"
	// DEV ...
	DEV = "/dev"
	// NOTIFYOOMONEY ...
	NOTIFYOOMONEY = "/notification/yoomoney"
)

// Handler -  Точка входа
func Handler(ctx context.Context, request *models.GatewayRequest) (*models.Response, error) {
	if request == nil {
		return nil, fmt.Errorf("nil request")
	}

	switch request.Path {
	case TG:
		return &models.Response{
			StatusCode: 200,
			Body:       tgForm(getParam(request, "id")),
		}, nil
	case DEV:
		return &models.Response{
			StatusCode: 200,
			Body:       "this DEV" + getParam(request, "id") + " " + request.Path,
		}, nil
	case NOTIFYOOMONEY:
		notification, err := notifications.NewYooMoney(*request)
		if err != nil {
			return nil, err
		}

		log.Println(notification, err)

		return &models.Response{
			StatusCode: 200,
			Body:       "OK",
		}, nil
	default:
		return &models.Response{
			StatusCode: 200,
			Body:       "null",
		}, nil
	}
}

// ViewData ...
type ViewData struct {
	ID     string
	Label  string
	Wallet string
}

func tgForm(ID string) string {
	data := ViewData{
		Label:  "TG-" + ID,
		ID:     ID,
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

func getParam(r *models.GatewayRequest, name string) string {
	values := r.Params[name]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
