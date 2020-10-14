package main

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
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

type TgUser struct {
	ID     string  `bson:"_id"`
	User   *User   `bson:"user"`
	Counts *Counts `bson:"counts"`
	Token  string  `bson:"token"`
}

type User struct {
	LanguageCode string `bson:"language_code"`
	Referral     string `bson:"referral"`
	Lang         string `bson:"lang"`
}

type Counts struct {
	CountPayments int     `bson:"count_payments"`
	SumSpent      float64 `bson:"sum_spent"`
	Referrals     int     `bson:"referrals"`
	Coins         int     `bson:"coins"`
	Rub           float64 `bson:"rub"`
}

type MongoDB struct {
	Payments *mongo.Collection
	Users    *mongo.Collection
}

func NewMongoDB() (MongoDB, error) {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	client, err := mongo.Connect(ctx, clientOptions)
	mdb := MongoDB{}
	if err != nil {
		return mdb, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return mdb, err
	}
	mdb.Payments = client.Database(os.Getenv("DB_NAME")).Collection("payments")
	mdb.Users = client.Database(os.Getenv("DB_NAME")).Collection("users")

	return mdb, err
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
	_, err = m.Payments.InsertOne(ctx, payment)
	return err
}

func (m *MongoDB) UpdateUser(user string, coins int, spent float64) (referral, token, lang string, err error) {
	var TgUser *TgUser
	opts := options.FindOneAndUpdate().SetProjection(bson.D{
		primitive.E{
			Key:   "token",
			Value: 1,
		}, primitive.E{
			Key:   "user.lang",
			Value: 1,
		}, primitive.E{
			Key:   "user.referral",
			Value: 1,
		}})
	filter := bson.D{primitive.E{Key: "_id", Value: user}}
	update := bson.D{
		primitive.E{
			Key: "$inc",
			Value: bson.D{
				primitive.E{
					Key:   "counts.coins",
					Value: coins,
				},
				primitive.E{
					Key:   "counts.count_payments",
					Value: 1,
				},
				primitive.E{
					Key:   "counts.sum_spent",
					Value: math.Ceil(spent*100) / 100,
				},
			},
		},
	}
	if err := m.Users.FindOneAndUpdate(ctx, filter, update, opts).Decode(&TgUser); err != nil {
		//if err == mongo.ErrNoDocuments {
		return "", "", "", err
	}
	if TgUser.Token == "" {
		TgUser.Token = os.Getenv("OLD_BOT_TOKEN")
	}
	return TgUser.User.Referral, TgUser.Token, TgUser.User.Lang, nil
}

func (m *MongoDB) UpdateReferral(referral string, summa float64) (token, lang string, err error) {
	var TgUser *TgUser
	filter := bson.D{primitive.E{Key: "_id", Value: referral}}
	update := bson.D{
		primitive.E{
			Key: "$inc",
			Value: bson.D{
				primitive.E{
					Key:   "counts.rub",
					Value: math.Ceil(summa*100) / 100,
				},
			},
		},
	}
	opts := options.FindOneAndUpdate().SetProjection(bson.D{
		primitive.E{
			Key:   "token",
			Value: 1,
		}, primitive.E{
			Key:   "user.lang",
			Value: 1,
		}})
	if err := m.Users.FindOneAndUpdate(ctx, filter, update, opts).Decode(&TgUser); err != nil {
		//if err == mongo.ErrNoDocuments {
		return "", "", err
	}
	if TgUser.Token == "" {
		TgUser.Token = os.Getenv("OLD_BOT_TOKEN")
	}
	return TgUser.Token, TgUser.User.Lang, nil
}
