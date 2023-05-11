package xstatus

import (
	"fmt"
)

// ErrNotFound
type ErrNotFound struct {
	Resource  string
	LookupKey string
	Message   string
}

func (err *ErrNotFound) Error() string {
	return fmt.Sprintf("not found: %v, %v, %v", err.Resource, err.LookupKey, err.Message)
}

func (err *ErrNotFound) Reason() string {
	return "Invalid Argument"
}

// ErrInvalidArgument
type ErrInvalidArgument struct {
	Argument string
	Message  string
}

func (err *ErrInvalidArgument) Error() string {
	return fmt.Sprintf("invalid argument: %v, %v", err.Argument, err.Message)
}

func (err *ErrInvalidArgument) Reason() string {
	return "Invalid Argument"
}

// ErrDeadlineExceeded
type ErrDeadlineExceeded struct {
	Message string
}

func (err *ErrDeadlineExceeded) Error() string {
	return fmt.Sprintf("deadline exceeded: %v", err.Message)
}

func (err *ErrDeadlineExceeded) Reason() string {
	return "Deadline Exceeded"
}

// ErrInternal
type ErrInternal struct {
	Message string
}

func (err *ErrInternal) Error() string {
	return fmt.Sprintf("internal error: %v", err.Message)
}

func (err *ErrInternal) Reason() string {
	return "Internal Error"
}
