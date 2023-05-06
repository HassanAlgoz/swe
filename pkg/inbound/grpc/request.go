package grpc

import (
	"fmt"

	"google.golang.org/grpc/metadata"
)

type Header string

const (
	HeaderUserId    Header = "x-user-id"
	HeaderRequestId Header = "x-request-id"
)

func GetUserId(md metadata.MD) (string, error) {
	vals := md.Get(string(HeaderUserId))
	if len(vals) == 0 || len(vals[0]) == 0 {
		return "", fmt.Errorf("missing header value: %s", HeaderUserId)
	}
	return vals[0], nil
}

func GetRequestId(md metadata.MD) (string, error) {
	vals := md.Get(string(HeaderRequestId))
	if len(vals) == 0 || len(vals[0]) == 0 {
		return "", fmt.Errorf("missing header value: %s", HeaderRequestId)
	}
	return vals[0], nil
}
