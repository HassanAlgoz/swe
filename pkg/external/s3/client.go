package s3

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/hassanalgoz/swe/pkg/infra/logger"
	"github.com/spf13/viper"
)

var log = logger.Get()

var (
	once     sync.Once
	instance s3iface.S3API
)

func Init() {
	var err error
	once.Do(func() {
		if viper.GetString("app.env") == "prod" {
			// Create the S3 client using the default session.
			sess := session.Must(
				session.NewSession(
					aws.NewConfig().
						WithMaxRetries(viper.GetInt("s3.client.max_retries")),
				))
			instance = s3.New(sess)
		}
	})
	if err != nil {
		log.Fatal().Msgf(`failed to initialize "%s" gRPC: %v`, "notify.client", err)
	}
}

func New() s3iface.S3API {
	Init()
	return instance
}
