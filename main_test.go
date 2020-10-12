package main

import (
	"encoding/json"
	"log"
	"testing"
	"time"
)



func TestUpdate_Validate(t *testing.T) {
	// Correct data
	body1 := string("{\"operation_id\": \"1234567\", \"notification_type\": \"p2p-incoming\",	\"datetime\":" +
		" \"2011-07-01T09:00:00.000+04:00\",	\"sha1_hash\": \"a2ee4a9195f4a90e893cff4f62eeba0b662321f9\",	\"sender\": " +
		"\"41001XXXXXXXX\",	\"currency\": \"643\",	\"amount\": 300.00,	\"label\": \"YM.label.12345\",	\"codepro\":" +
		" false,	\"unaccepted\": false}")
	// code_pro is true
	body2 := string("{\"operation_id\": \"1234567\", \"notification_type\": \"p2p-incoming\",	\"datetime\":" +
		" \"2011-07-01T09:00:00.000+04:00\",	\"sha1_hash\": \"a2ee4a9195f4a90e893cff4f62eeba0b662321f9\",	\"sender\": " +
		"\"41001XXXXXXXX\",	\"currency\": \"643\",	\"amount\": 300.00,	\"label\": \"YM.label.12345\",	\"codepro\":" +
		" true,	\"unaccepted\": false}")
	// Invalid sha1_hash
	body3 := string("{\"operation_id\": \"1234567\", \"notification_type\": \"p2p-incoming\",	\"datetime\":" +
		" \"2011-07-01T09:00:00.000+04:00\",	\"sha1_hash\": \"a2ee4a9195f4390e893cff4f62eeba0b662321f9\",	\"sender\": " +
		"\"41001XXXXXXXX\",	\"currency\": \"643\",	\"amount\": 300.00,	\"label\": \"YM.label.12345\",	\"codepro\":" +
		" false,	\"unaccepted\": false}")

	type fields struct {
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
	type args struct {
		notificationSecret string
	}

	update := new(fields)
	invalidCodeProUpdate := new(fields)
	invalidSha1HahsUpdate := new(fields)
	err := json.Unmarshal([]byte(body1), &update)
	err2 := json.Unmarshal([]byte(body2), &invalidCodeProUpdate)
	err3 := json.Unmarshal([]byte(body3), &invalidSha1HahsUpdate)
	if err != nil || err2 != nil || err3 != nil {
		log.Println(err)
	}


	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Valid",
			fields: *update,
			args: args{notificationSecret: "01234567890ABCDEF01234567890"},
			want: true,
			
		},
		{
			name: "inValidSecret",
			fields: *update,
			args: args{notificationSecret: "01234567890ABCDE301234567890"},
			want: false,

		},
		{
			name: "inValidCodePro",
			fields: *invalidCodeProUpdate,
			args: args{notificationSecret: "01234567890ABCDEF01234567890"},
			want: false,

		},
		{
			name: "invalidSha1Hahs",
			fields: *invalidSha1HahsUpdate,
			args: args{notificationSecret: "01234567890ABCDEF01234567890"},
			want: false,

		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Update{
				OperationId:      tt.fields.OperationId,
				NotificationType: tt.fields.NotificationType,
				Datetime:         tt.fields.Datetime,
				Sha1Hash:         tt.fields.Sha1Hash,
				Sender:           tt.fields.Sender,
				Currency:         tt.fields.Currency,
				Amount:           tt.fields.Amount,
				WithdrawAmount:   tt.fields.WithdrawAmount,
				Label:            tt.fields.Label,
				LastName:         tt.fields.LastName,
				FirstName:        tt.fields.FirstName,
				FathersName:      tt.fields.FathersName,
				Zip:              tt.fields.Zip,
				City:             tt.fields.City,
				Street:           tt.fields.Street,
				Building:         tt.fields.Building,
				Suite:            tt.fields.Suite,
				Flat:             tt.fields.Flat,
				Phone:            tt.fields.Phone,
				Email:            tt.fields.Email,
				TestNotification: tt.fields.TestNotification,
				CodePro:          tt.fields.CodePro,
				Unaccepted:       tt.fields.Unaccepted,
			}
			if got := u.Validate(tt.args.notificationSecret); got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
