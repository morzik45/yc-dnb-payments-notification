package main

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var ctx = context.TODO()

type PaymentInDB struct {
	Amount           float64   `bson:"amount" json:"amount"`
	WithdrawAmount   float64   `bson:"withdraw_amount" json:"withdraw_amount"`
	NotificationType string    `bson:"notification_type" json:"notification_type"`
	OperationId      string    `bson:"operation_id" json:"operation_id"`
	Currency         string    `bson:"currency" json:"currency"`
	UtcDatetime      time.Time `bson:"utc_datetime"`
	Datetime         time.Time `bson:"datetime" json:"datetime"`
	Sender           string    `bson:"sender" json:"sender"`
	CodePro          bool      `bson:"codepro" json:"code_pro"`
	Label            string    `bson:"label" json:"label"`
	Sha1Hash         string    `bson:"sha1_hash" json:"sha1_hash"`
	Unaccepted       bool      `bson:"unaccepted" json:"unaccepted"`
	PayDone          bool      `bson:"pay_done"`
}

type MongoDB struct {
	Payments *mongo.Collection
	Users    *mongo.Collection
}

func NewMongoDB() MongoDB {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	mdb := MongoDB{}

	mdb.Payments = client.Database("deepnudebot").Collection("payments")
	mdb.Users = client.Database("deepnudebot").Collection("users")

	return mdb
}

func (m *MongoDB) SaveInDB(u *Update) error {
	payment := new(PaymentInDB)
	j, err := json.Marshal(u)
	if err != nil {
		return err
	}
	err = json.Unmarshal(j, &payment)
	if err != nil {
		return err
	}
	payment.UtcDatetime = time.Now()
	payment.PayDone = true
	_, err = m.Payments.InsertOne(ctx, PaymentInDB{})
	return err
}

func (m MongoDB) GetFromDB() error {
	panic("implement me")
}

func (m MongoDB) UpdateUser() error {
	panic("implement me")
}

func (m MongoDB) GetUser() error {
	panic("implement me")
}
