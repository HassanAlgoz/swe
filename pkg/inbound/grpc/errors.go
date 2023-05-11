package grpc

import (
	"github.com/hassanalgoz/swe/pkg/xstatus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToStatusError(err error) error {
	if e, ok := err.(*xstatus.ErrNotFound); ok {
		return status.Error(codes.NotFound, e.Error())
	} else if e, ok := err.(*xstatus.ErrInvalidArgument); ok {
		return status.Error(codes.InvalidArgument, e.Error())
	}
	return status.Error(codes.Internal, err.Error())
}
