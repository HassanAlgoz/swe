package port

import (
	"context"

	"github.com/hassanalgoz/swe/internal/services/lms/controller"
	"github.com/hassanalgoz/swe/pkg/entities"
	"github.com/hassanalgoz/swe/pkg/infra/logger"
	inbound "github.com/hassanalgoz/swe/pkg/ports/inbound/grpc"
	pb "github.com/hassanalgoz/swe/ports/services/lms"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type service struct {
	pb.UnimplementedLMSServer
}

var log = logger.Singleton()
var ctrl = controller.Singleton()

func Register(registrar grpc.ServiceRegistrar) {
	pb.RegisterLMSServer(registrar, service{})
}

func (s service) CreateCourse(ctx context.Context, req *pb.CreateCourseRequest) (*pb.CreateCourseResponse, error) {
	log.Debug().Msgf("[CreateCourse] start")

	// Headers
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "missing request metadata")
	}
	userId, err := inbound.GetUserId(md)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	reqId, err := inbound.GetRequestId(md)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	log.Info().
		Str("userId", userId).
		Str("reqId", reqId).Send()

	id, err := ctrl.CreateCourse(ctx, entities.Course{
		Name:        req.GetName(),
		Description: req.GetDescription(),
	})
	if err != nil {
		inbound.ToStatusError(err)
	}
	log.Debug().Msgf("[CreateCourse] end")
	return &pb.CreateCourseResponse{
		Id: id.String(),
	}, nil
}
