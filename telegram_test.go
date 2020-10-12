package main

import "testing"

func TestTelegram_SendNotification(t1 *testing.T) {
	type args struct {
		user  string
		token string
		text  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				user:  "<PAST>",
				token: "<PAST>",
				text: "Поступил платёж на сумму: <b>399₽</b>\nНа твой счёт зачислено <b>100</b> монет (из них <b>20</b> это бонус).\n" +
					"Проверить баланс ты можешь командой /info. Спасибо что помогаете развитию бота.",
			},
			wantErr: false,
		},
		{
			name: "not ok",
			args: args{
				user:  "<PAST>",
				token: "<PAST>",
				text: "Поступил платёж на сумму: %.2f₽\nНа твой счёт зачислено %d монет (из них %d это бонус).\n" +
					"Проверить баланс ты можешь командой /info. Спасибо что помогаете развитию бота.",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := Telegram{}
			if err := t.SendNotification(tt.args.user, tt.args.token, tt.args.text); (err != nil) != tt.wantErr {
				t1.Errorf("SendNotification() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
