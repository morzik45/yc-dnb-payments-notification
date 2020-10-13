package main

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type RequestBody struct {
	Body string `json:"body"`
}

type Response struct {
	StatusCode int `json:"statusCode"`
}

type DB interface {
	SaveInDB(u *Update) error
	UpdateUser(user string, coins int, spent float64) (referral, token, lang string, err error)
	UpdateReferral(referral string, summa float64) (token, lang string, err error)
}

type Notification interface {
	SendNotification(user, token, text string) error
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

func (u *Update) Bonus() (coins, bonus int) {
	coinCost := 4
	coins = int(math.Ceil(math.Ceil(u.Amount) / float64(coinCost)))
	if u.Amount > 975 {
		bonus = 500
	} else if u.Amount > 490 {
		bonus = 112
	} else if u.Amount > 390 {
		bonus = 80
	} else if u.Amount > 290 {
		bonus = 50
	} else if u.Amount > 190 {
		bonus = 30
	} else if u.Amount > 90 {
		bonus = 12
	}
	return
}

func (u *Update) Processes(db DB, msg Notification) error {

	// Сохраняем факт поступления платежа в базу
	if err := db.SaveInDB(u); err != nil {
		SaveError("errors", fmt.Sprintf("func: db.SaveInDB\nerror: %s\nUpdate: %v", err, &u))
	}

	// Обновляем баланс пользователя
	coins, bonus := u.Bonus()
	referral, uToken, uLang, err := db.UpdateUser(u.Label, coins+bonus, u.Amount)
	if err != nil {
		return err
	}

	// Уведомляем пользователя
	var uText string
	if uLang == "ru" {
		uText = fmt.Sprintf("Поступил платёж на сумму: <b>%.2f₽</b>\nНа твой счёт зачислено <b>%d</b> монет (из них <b>%d</b> это бонус).\n"+
			"Проверить баланс ты можешь командой /info. Спасибо что помогаете развитию бота.", u.Amount, coins+bonus, bonus)
	} else {
		uText = fmt.Sprintf("Received a payment in the amount of <b>%.2f₽</b>\n<b>%d</b> coins were credited to your account"+
			" (<b>%d</b> of them are a bonus).\nYou can check the balance with the /info command. Thank you for helping"+
			" the bot develop.", u.Amount, coins+bonus, bonus)
	}
	if err := msg.SendNotification(u.Label, uToken, uText); err != nil {
		SaveError("errors", fmt.Sprintf("func: msg.SendNotification\nerror: %s\nUpdate: %v", err, &u))
	}

	// Если есть реферрал
	if referral != "" {

		// Обновляем баланс реферрала
		summa := math.Ceil((u.Amount*100/50)*100) / 100
		rToken, rLang, err := db.UpdateReferral(referral, summa)
		if err != nil {
			SaveError("errors", fmt.Sprintf("func: db.UpdateReferral\nerror: %s\nUpdate: %v", err, &u))
		}

		// Уведомляем реферрала
		var rText string
		if rLang == "ru" {
			rText = fmt.Sprintf("+ <b>%.2f ₽</b> 💰\nПодробнее /info", summa)
		} else {
			rText = fmt.Sprintf("+ <b>%.2f ₽</b> 💰\nDetails /info", summa)
		}
		if err := msg.SendNotification(referral, rToken, rText); err != nil {
			SaveError("errors", fmt.Sprintf("func: msg.SendNotification2\nerror: %s\nUpdate: %v", err, &u))
		}
	}

	// Уведомляем о платеже админов
	if err := msg.SendNotification(
		os.Getenv("PAYMENTS_CHAT"),
		"ADMIN_BOT_TOKEN",
		fmt.Sprintf("Новый плтёж на сумму <b>%.2f</b> от пользователя <i>%s</i> (<code>%s</code>)\n"+
			"Реферрал: <i>%s</i>", u.Amount, u.Label, u.OperationId, referral),
	); err != nil {
		SaveError("errors", fmt.Sprintf("func: msg.SendNotification3\nerror: %s\nUpdate: %v", err, &u))
	}

	return nil
}

func toJSON(m interface{}) string {
	js, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	s := strings.ReplaceAll(string(js), ",", ", ")
	return strings.ReplaceAll(s, "%3A", ":")
}

func Handler(_ context.Context, request RequestBody) (*Response, error) {
	bytesBody, err := base64.StdEncoding.DecodeString(request.Body) // Converting data
	if err != nil {
		SaveError("errors", fmt.Sprintf("Failed to Decode secret\nfunc: base64.StdEncoding.DecodeString\nerror: %s\nUpdate: %s", err, request.Body))
	}
	a, err := url.ParseQuery(string(bytesBody))
	if err != nil {
		SaveError("errors", fmt.Sprintf("func: url.ParseQuery\nerror: %s\nUpdate: %s", err, request.Body))
	}
	decoder := json.NewDecoder(bytes.NewReader([]byte(toJSON(a))))
	update := new(Update)
	err = decoder.Decode(&update)
	if err != nil {
		SaveError("errors", fmt.Sprintf("func: Handler_json.Unmarshal\nerror: %s\nUpdate: %s", err, string(bytesBody)))
	} else if update.Validate(os.Getenv("YM_SECRET")) {
		//	custom logic
		mdb, err := NewMongoDB()
		if err != nil {
			SaveError("errors", fmt.Sprintf("func: NewMongoDB\nerror: %s\nUpdate: %s", err, string(bytesBody)))
		} else if err = update.Processes(&mdb, Telegram{}); err != nil {
			SaveError("errors", fmt.Sprintf("func: update.Processes\nerror: %s\nUpdate: %s", err, string(bytesBody)))
		}
	}
	return &Response{StatusCode: http.StatusOK}, nil
}
