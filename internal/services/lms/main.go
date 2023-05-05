package main

import (
	"net"

	"github.com/hassanalgoz/swe/internal/services/lms/view"
	"github.com/hassanalgoz/swe/pkg/infra/logger"
	"google.golang.org/grpc"
)

var log = logger.Singleton()

func main() {
	// Initialize
	server := grpc.NewServer()
	view.Register(server)

	// Listen and serve
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}
	if err := server.Serve(lis); err != nil {
		log.Fatal().Msgf("failed to serve: %v", err)
	}
}
