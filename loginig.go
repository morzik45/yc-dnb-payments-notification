package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"time"
)

func SaveError(level, myErr string) {
	bucket := os.Getenv("DB_NAME")
	filename := time.Now().Format(time.RFC3339) + ".txt"

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("ru-central1"),
		Endpoint: aws.String("https://storage.yandexcloud.net"),
	},
	)
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fmt.Sprintf("logs/payments/%s/%s", level, filename)),
		Body:   bytes.NewReader([]byte(myErr)),
	})
	if err != nil {
		fmt.Errorf("unable to upload %q to %q, %v", filename, bucket, err)
	}
	fmt.Printf("Successfully uploaded %q to %q\n", filename, bucket)
}