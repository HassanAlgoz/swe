package main

import (
	"net"

	"github.com/hassanalgoz/swe/internal/services/lms/controller"
	"github.com/hassanalgoz/swe/internal/services/lms/port"
	"github.com/hassanalgoz/swe/internal/services/lms/store"
	"github.com/hassanalgoz/swe/pkg/external/s3"
	"github.com/hassanalgoz/swe/pkg/infra/logger"
	"github.com/hassanalgoz/swe/pkg/services/adapters/notify"
	"google.golang.org/grpc"
)

var log = logger.Get()

func main() {
	ctrl := controller.New(
		store.New("lms"),
		notify.New(),
		s3.New(),
	)

	// Initialize
	server := grpc.NewServer()
	port.Register(server, ctrl)

	// Listen and serve
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}
	if err := server.Serve(lis); err != nil {
		log.Fatal().Msgf("failed to serve: %v", err)
	}
}
