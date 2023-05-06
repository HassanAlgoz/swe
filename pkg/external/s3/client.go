package s3

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/spf13/viper"
)

var (
	instance s3iface.S3API
	once     sync.Once
)

func Singleton() s3iface.S3API {
	switch viper.GetString("app.env") {
	case "prod":
		once.Do(func() {
			// Create the S3 client using the default session.
			sess := session.Must(session.NewSession(aws.NewConfig().WithMaxRetries(viper.GetInt("s3.client.max_retries"))))
			instance = s3.New(sess)
		})
	default:
		once.Do(func() {
			instance = newMockedS3Client()
		})
	}
	return instance
}
