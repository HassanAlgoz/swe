package notify

import (
	"fmt"
	"sync"

	"github.com/hassanalgoz/swe/pkg/infra/logger"
	port "github.com/hassanalgoz/swe/pkg/services/ports/notify"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var log = logger.Get()

var (
	once sync.Once
	conn *grpc.ClientConn
)

type Adapter struct {
	port.NotificationsClient
	// namespace string
}

func Init() {
	var err error
	once.Do(func() {
		target := fmt.Sprintf("%s:%s",
			viper.GetString("services.lms.host"),
			viper.GetString("services.lms.port"))
		conn, err = grpc.Dial(target, grpc.WithInsecure())
	})
	if err != nil {
		log.Fatal().Msgf(`failed to initialize "%s" gRPC: %v`, "notify.client", err)
	}
}

func New() port.NotificationsClient {
	Init()
	client := &Adapter{
		NotificationsClient: port.NewNotificationsClient(conn),
	}
	return client
}

// TODO: this shall be used in integration test, but I am not there yet.
// func (a *Adapter) appendDefaultMetadata(ctx context.Context) {
// 	if viper.GetString("env") == "test" {
// 		metadata.AppendToOutgoingContext(ctx,
// 			"namespace", a.namespace,
// 		)
// 	}
// }
