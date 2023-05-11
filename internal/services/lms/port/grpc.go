package port

import (
	"context"

	"github.com/hassanalgoz/swe/internal/services/lms/controller"
	StorePort "github.com/hassanalgoz/swe/internal/services/lms/store/port"
	inbound "github.com/hassanalgoz/swe/pkg/inbound/grpc"
	"github.com/hassanalgoz/swe/pkg/infra/logger"
	pb "github.com/hassanalgoz/swe/pkg/services/ports/lms"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type service struct {
	pb.UnimplementedLMSServer
	controller *controller.Controller
}

var log = logger.Get()

func Register(registrar grpc.ServiceRegistrar, ctrl *controller.Controller) {
	pb.RegisterLMSServer(registrar, service{
		controller: ctrl,
	})
}

func (s service) CreateCourse(ctx context.Context, req *pb.CoursePut) (*pb.Course, error) {
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

	// Log
	log.Info().
		Str("userId", userId).
		Str("reqId", reqId).Send()

	// call and return
	course, err := s.controller.CreateCourse(ctx, StorePort.Course{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		inbound.ToStatusError(err)
	}

	log.Debug().Msgf("[CreateCourse] end")
	return &pb.Course{
		Id:          course.ID.String(),
		Code:        course.Code,
		Name:        course.Name,
		Description: course.Description,
	}, nil
}
