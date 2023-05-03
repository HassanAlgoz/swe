package http

import (
	"encoding/json"
	"net/http"

	"github.com/hassanalgoz/swe/pkg/entities"
)

type LocationTypeEnum string

const (
	LocationTypeHeader    LocationTypeEnum = "header"
	LocationTypeParameter LocationTypeEnum = "parameter"
)

// Response conforms to: https://google.github.io/styleguide/jsoncstyleguide.xml
type Response struct {
	Data  any   `json:"data,omitempty"`
	Error Error `json:"error,omitempty"`
}

type Error struct {
	// This property value will usually represent the HTTP response code.
	// If there are multiple errors, code will be the error code for the first error.
	Code int `json:"code"`

	// A human readable message providing more details about the error.
	// If there are multiple errors, message will be the message for the first error.
	Message string `json:"message"`

	// Container for any additional information regarding the error.
	// If the service returns multiple errors, each element in the errors array represents a different error.
	Errors []ErrorItem `json:"errors"`
}

type ErrorItem struct {
	// A human readable message providing more details about the error.
	// If there is only one error, this field will match error.message.
	Message string `json:"message"`

	// Unique identifier for this error.
	// Different from the error.code property in that this is not an http response code.
	Reason string `json:"reason"`

	// if LocationType = "header" then it may be: "Authorization
	// if LocationType = "parameter" then it may be: "orderId"
	Location     string           `json:"location"`
	LocationType LocationTypeEnum `json:"location_type"`
}

// Ok 200
func Ok(w http.ResponseWriter, data any) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Data: data,
	})
}

// ErrNotFound 404
func ErrNotFound(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{
		Error: Error{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		},
	})
}

// ErrInternal 500
func ErrInternal(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(Response{
		Error: Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	})
}

// ErrInvalidArgument 400
func ErrInvalidArgument(w http.ResponseWriter, err *entities.ErrInvalidArgument) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(Response{
		Error: Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Errors: []ErrorItem{
				{
					Message:      err.Error(),
					Reason:       err.Reason(),
					LocationType: LocationTypeParameter,
					Location:     err.Argument,
				},
			},
		},
	})
}

// ErrInvalidState 200
func ErrInvalidState(w http.ResponseWriter, err *entities.ErrInvalidState) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Error: Error{
			Code:    http.StatusOK,
			Message: err.Error(),
			Errors: []ErrorItem{
				{
					Message:      err.Error(),
					Reason:       err.Reason(),
					LocationType: LocationTypeParameter,
					Location:     err.RelatedArgument,
				},
			},
		},
	})
}
