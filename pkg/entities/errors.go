package entities

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
