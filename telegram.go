package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TGResponse struct {
	Ok bool `json:"ok"`
}

type Telegram struct{}

func (t Telegram) SendNotification(user, token, text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	requestBody, err := json.Marshal(map[string]string{
		"chat_id":    user,
		"text":       text,
		"parse_mode": "HTML",
	})
	if err != nil {
		SaveError(fmt.Sprintf("ошибка при кодировании json при отправке уведомления: %s", err))
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		SaveError(fmt.Sprintf("ошибка при отправке запроса в TG: %s", err))
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			SaveError(fmt.Sprintf("ошибка при закрытии тела ответа от тг: %s", err))
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		SaveError(fmt.Sprintf("ошибка при расшифровке ответа от TG: %s", err))
		return err
	}
	tr := new(TGResponse)
	err = json.Unmarshal(body, &tr)
	if !tr.Ok {
		SaveError(fmt.Sprintf("TG вернул ошибку: %s", string(body)))
		return errors.New(string(body))
	}
	return nil

}
