package main

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"
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
func Handler(ctx context.Context, request *GatewayRequest) (*Response, error) {
	if request == nil {
		return nil, fmt.Errorf("nil request")
	}

	switch request.Path {
	case TG:
		return &Response{
			StatusCode: 200,
			Body:       tgForm(getParam(request, "id")),
		}, nil
	case DEV:
		return &Response{
			StatusCode: 200,
			Body:       "this DEV" + getParam(request, "id") + " " + request.Path,
		}, nil
	case NOTIFYOOMONEY:
		return YooMoneyNotifFunc(request)
	default:
		return &Response{
			StatusCode: 200,
			Body:       "null",
		}, nil
	}
}

// YooMoneyNotifFunc ...
func YooMoneyNotifFunc(request *GatewayRequest) (*Response, error) {
	notif, err := NewUpdateYooMoney(*request)
	if err != nil {
		return nil, err
	}

	log.Println(notif)

	return &Response{
		StatusCode: 200,
		Body:       "OK",
	}, nil
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

// NotifYooMoney ...
type NotifYooMoney struct {
	OperationID      string    `json:"operation_id"`
	NotificationType string    `json:"notification_type"`
	Datetime         time.Time `json:"datetime"`
	Sha1Hash         string    `json:"sha1_hash"`
	Sender           string    `json:"sender"`
	Currency         string    `json:"currency"`
	Amount           float64   `json:"amount"`
	WithdrawAmount   float64   `json:"withdraw_amount"`
	Label            string    `json:"label"`
	LastName         string    `json:"lastname"`
	FirstName        string    `json:"firstname"`
	FathersName      string    `json:"fathersname"`
	Zip              string    `json:"zip"`
	City             string    `json:"city"`
	Street           string    `json:"street"`
	Building         string    `json:"building"`
	Suite            string    `json:"suite"`
	Flat             string    `json:"flat"`
	Phone            string    `json:"phone"`
	Email            string    `json:"email"`
	TestNotification bool      `json:"test_notification"`
	CodePro          bool      `json:"codepro"`
	Unaccepted       bool      `json:"unaccepted"`
}

// ValidateYooMoney ...
func (u *NotifYooMoney) ValidateYooMoney(notificationSecret string) bool {
	s := fmt.Sprintf("%s&%s&%.2f&%s&%s&%s&%t&%s&%s",
		u.NotificationType,
		u.OperationID,
		u.Amount,
		u.Currency,
		u.Datetime.Format("2006-01-02T03:04:05Z"),
		u.Sender,
		u.CodePro,
		notificationSecret,
		u.Label,
	)
	h := sha1.New()
	h.Write([]byte(s))
	mySha1Hash := hex.EncodeToString(h.Sum(nil))
	if mySha1Hash != u.Sha1Hash || u.CodePro || u.Unaccepted {
		return false
	}
	return true
}

// NewUpdateYooMoney ...
func NewUpdateYooMoney(request GatewayRequest) (*NotifYooMoney, error) {
	update := new(NotifYooMoney)
	bytesBody, err := base64.StdEncoding.DecodeString(request.Body) // Converting data
	if err != nil {
		return update, err
	}
	a, err := url.ParseQuery(string(bytesBody))
	if err != nil {
		return update, err
	}
	update.OperationID = a.Get("operation_id")
	update.NotificationType = a.Get("notification_type")
	update.Datetime, err = time.Parse(time.RFC3339, a.Get("datetime"))
	if err != nil {
		return update, err
	}
	update.Sha1Hash = a.Get("sha1_hash")
	update.Sender = a.Get("sender")
	update.Currency = a.Get("currency")
	update.Amount, err = strconv.ParseFloat(a.Get("amount"), 64)
	if err != nil {
		return update, err
	}
	update.WithdrawAmount, err = strconv.ParseFloat(a.Get("withdraw_amount"), 64)
	if err != nil && a.Get("withdraw_amount") != "" {
		return update, err
	}
	update.Label = a.Get("label")
	update.LastName = a.Get("lastname")
	update.FirstName = a.Get("firstname")
	update.FathersName = a.Get("fathersname")
	update.Zip = a.Get("zip")
	update.City = a.Get("city")
	update.Street = a.Get("street")
	update.Building = a.Get("building")
	update.Suite = a.Get("suite")
	update.Flat = a.Get("flat")
	update.Phone = a.Get("phone")
	update.Email = a.Get("email")
	update.TestNotification, err = strconv.ParseBool(a.Get("test_notification"))
	if err != nil && a.Get("withdraw_amount") != "" {
		return update, err
	}
	update.CodePro, err = strconv.ParseBool(a.Get("codepro"))
	if err != nil && a.Get("withdraw_amount") != "" {
		return update, err
	}
	update.Unaccepted, err = strconv.ParseBool(a.Get("unaccepted"))
	if err != nil && a.Get("withdraw_amount") != "" {
		return update, err
	}
	return update, nil
}

func getParam(r *GatewayRequest, name string) string {
	values := r.Params[name]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
