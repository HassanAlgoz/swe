package lms

import (
	"fmt"
	"sync"
	"time"

	"github.com/hassanalgoz/swe/pkg/infra/logger"
	port "github.com/hassanalgoz/swe/pkg/services/ports/lms"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var log = logger.Get()

var (
	once sync.Once
	conn *grpc.ClientConn
)

type Adapter struct {
	port.LMSClient
}

func Init() {
	var err error
	once.Do(func() {
		target := fmt.Sprintf("%s:%s",
			viper.GetString("services.default.host"),
			viper.GetString("services.default.port"))
		conn, err = grpc.Dial(
			target,
			grpc.WithInsecure(),
			grpc.WithTimeout(time.Duration(viper.GetInt("grpc.client.timeout"))*time.Second),
		)
	})
	if err != nil {
		log.Fatal().Msgf(`failed to initialize "%s" gRPC: %v`, "notify.client", err)
	}
}

func New() port.LMSClient {
	Init()
	client := &Adapter{
		LMSClient: port.NewLMSClient(conn),
	}
	return client
}
