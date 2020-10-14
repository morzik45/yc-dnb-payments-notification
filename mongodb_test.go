package main

import (
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"testing"
	"time"
)

var mdb MongoDB

func TestNewMongoDB(t *testing.T) {
	tests := []struct {
		name    string
		want    MongoDB
		wantErr bool
	}{
		{
			name:    "get DB",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			mdb, err = NewMongoDB()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMongoDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMongoDB_SaveInDB(t *testing.T) {
	type fields struct {
		Payments *mongo.Collection
		Users    *mongo.Collection
	}
	type args struct {
		u *Update
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create test payment",
			fields: fields{
				Payments: mdb.Payments,
				Users:    mdb.Users,
			},
			args: args{&Update{
				OperationId:      "a1b2c3",
				NotificationType: "test-notification",
				Datetime:         time.Now(),
				Sha1Hash:         "3df21g3df2g16d5h1gf3h12g",
				Sender:           "987654321",
				Currency:         "643",
				Amount:           99.99,
				WithdrawAmount:   0,
				Label:            "123456789",
				TestNotification: true,
				CodePro:          false,
				Unaccepted:       false,
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MongoDB{
				Payments: tt.fields.Payments,
				Users:    tt.fields.Users,
			}
			if err := m.SaveInDB(tt.args.u); (err != nil) != tt.wantErr {
				t.Errorf("SaveInDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMongoDB_UpdateReferral(t *testing.T) {
	type fields struct {
		Payments *mongo.Collection
		Users    *mongo.Collection
	}
	type args struct {
		referral string
		summa    float64
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantToken string
		wantLang  string
		wantErr   bool
	}{
		{
			name: "update me as referral",
			fields: fields{
				Payments: mdb.Payments,
				Users:    mdb.Users,
			},
			args: args{
				referral: os.Getenv("ADMIN_ID"),
				summa:    10.01,
			},
			wantToken: os.Getenv("OLD_BOT_TOKEN"),
			wantLang:  "ru",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MongoDB{
				Payments: tt.fields.Payments,
				Users:    tt.fields.Users,
			}
			gotToken, gotLang, err := m.UpdateReferral(tt.args.referral, tt.args.summa)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateReferral() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotToken != tt.wantToken {
				t.Errorf("UpdateReferral() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
			if gotLang != tt.wantLang {
				t.Errorf("UpdateReferral() gotLang = %v, want %v", gotLang, tt.wantLang)
			}
		})
	}
}

func TestMongoDB_UpdateUser(t *testing.T) {
	type fields struct {
		Payments *mongo.Collection
		Users    *mongo.Collection
	}
	type args struct {
		user  string
		coins int
		spent float64
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantReferral string
		wantToken    string
		wantLang     string
		wantErr      bool
	}{
		{
			name: "get me",
			fields: fields{
				Payments: mdb.Payments,
				Users:    mdb.Users,
			},
			args: args{
				user:  os.Getenv("ADMIN_ID"),
				coins: 3,
				spent: 100.34,
			},
			wantReferral: "885146682",
			wantToken:    os.Getenv("OLD_BOT_TOKEN"),
			wantLang:     "ru",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MongoDB{
				Payments: tt.fields.Payments,
				Users:    tt.fields.Users,
			}
			gotReferral, gotToken, gotLang, err := m.UpdateUser(tt.args.user, tt.args.coins, tt.args.spent)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotReferral != tt.wantReferral {
				t.Errorf("UpdateUser() gotReferral = %v, want %v", gotReferral, tt.wantReferral)
			}
			if gotToken != tt.wantToken {
				t.Errorf("UpdateUser() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
			if gotLang != tt.wantLang {
				t.Errorf("UpdateUser() gotLang = %v, want %v", gotLang, tt.wantLang)
			}
		})
	}
}
