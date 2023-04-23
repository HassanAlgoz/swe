package common

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("not found")
)

// ErrInvalidArgument indicates an error with passed inputs
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

// ErrInvalidState indicates an error with persisted state related to inputs
type ErrInvalidState struct {
	RelatedArgument string
	Message         string
}

func (err *ErrInvalidState) Error() string {
	return fmt.Sprintf("invalid state related to: %v, %v", err.RelatedArgument, err.Message)
}

func (err *ErrInvalidState) Reason() string {
	return "Invalid State"
}
