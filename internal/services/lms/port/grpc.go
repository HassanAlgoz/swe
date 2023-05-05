package port

import (
	"context"
	"errors"

	"github.com/hassanalgoz/swe/internal/services/lms/controller"
	"github.com/hassanalgoz/swe/pkg/entities"
	"github.com/hassanalgoz/swe/pkg/infra/logger"
	pb "github.com/hassanalgoz/swe/ports/services/lms"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
	// TODO: extract and use request header data (authorization)
	id, err := ctrl.CreateCourse(ctx, entities.Course{
		Name:        req.GetName(),
		Description: req.GetDescription(),
	})
	if err != nil {
		// TODO: fix status code (maybe define general application codes in a proto first)
		if errors.Is(err, entities.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		} else if e, ok := err.(*entities.ErrInvalidArgument); ok {
			return nil, status.Error(codes.InvalidArgument, e.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	log.Debug().Msgf("[CreateCourse] end")
	return &pb.CreateCourseResponse{
		Id: id.String(),
	}, nil
}
