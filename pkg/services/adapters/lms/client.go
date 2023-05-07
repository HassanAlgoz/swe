package lms

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/pkg/infra/logger"
	port "github.com/hassanalgoz/swe/pkg/services/ports/lms"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type Adapter struct {
	port port.LMSClient
}

var (
	once     sync.Once
	instance *Adapter
)

func Singleton() *Adapter {
	var err error
	log := logger.Get()

	once.Do(func() {
		var conn *grpc.ClientConn
		conn, err = grpc.Dial(fmt.Sprintf("%s:%s", viper.GetString("grpc.search.host"), viper.GetString("grpc.search.port")), grpc.WithInsecure())
		if err != nil {
			return
		}
		instance = &Adapter{
			port: port.NewLMSClient(conn),
		}
	})
	if err != nil {
		log.Fatal().Msgf("failed to instantiate gRPC client adapter: %v", err)
	}
	return instance
}

func (a *Adapter) CreateCourse(ctx context.Context, req *port.CreateCourseRequest) (*uuid.UUID, error) {
	resp, err := a.port.CreateCourse(ctx, req)
	if err != nil {
		return nil, err
	}
	id := uuid.MustParse(resp.Id)
	return &id, nil
}
