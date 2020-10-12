package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
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

type DB interface {
	SaveInDB(u *Update) error
	UpdateUser(user string, coins, bonus int) (referral, token, lang string, err error)
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
		fmt.Errorf("ошибка при сохранении платежа в базу: %s", err)
	}

	// Обновляем баланс пользователя
	coins, bonus := u.Bonus()
	referral, uToken, uLang, err := db.UpdateUser(u.Label, coins, bonus)
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
		fmt.Errorf("ошибка при уведомлении пользователя: %s", err)
	}

	// Если есть реферрал
	if referral != "" {

		// Обновляем баланс реферрала
		var summa float64 = math.Ceil((u.Amount*100/50)*100) / 100
		rToken, rLang, err := db.UpdateReferral(referral, summa)
		if err != nil {
			fmt.Errorf("ошибка при изменении аккаунта реферрала: %s", err)
		}

		// Уведомляем реферрала
		var rText string
		if rLang == "ru" {
			rText = fmt.Sprintf("+ <b>%.2f ₽</b> 💰\nПодробнее /info", summa)
		} else {
			rText = fmt.Sprintf("+ <b>%.2f ₽</b> 💰\nDetails /info", summa)
		}
		if err := msg.SendNotification(referral, rToken, rText); err != nil {
			fmt.Errorf("ошибка при уведомлении реферрала: %s", err)
		}
	}

	// Уведомляем о платеже админов
	if err := msg.SendNotification(
		os.Getenv("PAYMENTS_CHAT"),
		"ADMIN_BOT_TOKEN",
		fmt.Sprintf("Новый плтёж на сумму <b>%.2f</b> от пользователя <i>%s</i> (<code>%s</code>)\n"+
			"Реферрал: <i>%s</i>", u.Amount, u.Label, u.OperationId, referral),
	); err != nil {
		fmt.Errorf("ошибка при уведомлении админов: %s", err)
	}

	return nil
}

func Handler(_ context.Context, request RequestBody) (*Response, error) {
	update := new(Update)
	err := json.Unmarshal([]byte(request.Body), &update)
	if err != nil {
		log.Println(err)
	}
	if update.Validate(os.Getenv("YM_SECRET")) {
		//	custom logic
		mdb := NewMongoDB()
		if err = update.Processes(&mdb, Telegram{}); err != nil {

		}

	}
	return &Response{StatusCode: http.StatusOK}, nil
}
