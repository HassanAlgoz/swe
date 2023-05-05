package notify

import (
	"context"
	"fmt"
	"sync"

	"github.com/hassanalgoz/swe/pkg/infra/logger"
	port "github.com/hassanalgoz/swe/ports/services/notify"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type Adapter struct {
	port port.NotificationsClient
}

var (
	once     sync.Once
	instance *Adapter
)

func Singleton() *Adapter {
	var err error
	log := logger.Singleton()

	once.Do(func() {
		var conn *grpc.ClientConn
		conn, err = grpc.Dial(fmt.Sprintf("%s:%s", viper.GetString("grpc.search.host"), viper.GetString("grpc.search.port")), grpc.WithInsecure())
		if err != nil {
			return
		}
		instance = &Adapter{
			port: port.NewNotificationsClient(conn),
		}
	})
	if err != nil {
		log.Fatal().Msgf("failed to instantiate gRPC client adapter: %v", err)
	}
	return instance
}

func (a *Adapter) SendNotification(ctx context.Context, req *port.NotificationRequest) error {
	_, err := a.port.SendNotification(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
