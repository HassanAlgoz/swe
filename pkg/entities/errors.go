package entities

import (
	"fmt"
)

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
