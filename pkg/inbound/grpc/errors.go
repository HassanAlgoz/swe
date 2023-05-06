package grpc

import (
	"errors"

	"github.com/hassanalgoz/swe/pkg/entities"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToStatusError(err error) error {
	if errors.Is(err, entities.ErrNotFound) {
	} else if e, ok := err.(*entities.ErrInvalidArgument); ok {
		return status.Error(codes.InvalidArgument, e.Error())
	}
	return status.Error(codes.Internal, err.Error())
}
