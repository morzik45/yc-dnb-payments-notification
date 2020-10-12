package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type RequestBody struct {
	Body string `json:"body"`
}

type Response struct {
	StatusCode int `json:"statusCode"`
}

type Update struct {
	OperationId      string    `json:"operation_id"`
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

func (u *Update) Validate(notificationSecret string) bool {
	s := fmt.Sprintf("%s&%s&%.2f&%s&%s&%s&%t&%s&%s",
		u.NotificationType,
		u.OperationId,
		u.Amount,
		u.Currency,
		u.Datetime.Format("2006-01-02T03:04:05.000-07:00"),
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

type DB interface {
	SaveInDB(u *Update) error
	GetFromDB() error
	UpdateUser() error
	GetUser() error
}

type Notification interface {
	SendNotification(text, user string) error
}

func Handler(_ context.Context, request RequestBody) (*Response, error) {
	update := new(Update)
	err := json.Unmarshal([]byte(request.Body), &update)
	if err != nil {
		log.Println(err)
	}
	if update.Validate(os.Getenv("YM_SECRET")) {
	//	You custom logic


	}
	return &Response{StatusCode: http.StatusOK}, nil
}
