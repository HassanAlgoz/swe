package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/internal/actions"
	"github.com/hassanalgoz/swe/internal/entities"
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

func registerHandlers(ctx context.Context, mux *http.ServeMux, act actions.Actions) {
	mux.HandleFunc("/actions:transfer-money", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(Response{
				Error: Error{
					Code:    http.StatusMethodNotAllowed,
					Message: "Method Not Allowed",
				},
			})
			return
		}

		// Read request headers, parse, validate, ...etc.
		userID := r.Header.Get("x-user-id")
		// Parse, validate.....
		if userID == "" {
			userID = "default-user-id"
		}

		// Apply rules based on headers...
		// ...
		// ...

		// Parse body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{
				Error: Error{
					Code:    http.StatusBadRequest,
					Message: "Invalid request body",
				},
			})
			return
		}

		// Parse body.json
		var fields struct {
			From   string `json:"from"`
			To     string `json:"to"`
			Amount int64  `json:"amount"`
		}
		err = json.Unmarshal(body, &fields)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{
				Error: Error{
					Code:    http.StatusBadRequest,
					Message: "Invalid json",
				},
			})
			return
		}

		// Parse body.json.fields
		// note: errors are appended
		var fieldsErrors []ErrorItem
		from, err := uuid.Parse(fields.From)
		if err != nil {
			fieldsErrors = append(fieldsErrors, ErrorItem{
				LocationType: "parameter",
				Location:     "from",
				Message:      "invalid uuid",
				Reason:       err.Error(),
			})
		}
		to, err := uuid.Parse(fields.To)
		if err != nil {
			fieldsErrors = append(fieldsErrors, ErrorItem{
				LocationType: "parameter",
				Location:     "to",
				Message:      "invalid uuid",
				Reason:       err.Error(),
			})
		}

		amount := fields.Amount

		// Return parsing errors (if any)
		if len(fieldsErrors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{
				Error: Error{
					Code:    http.StatusBadRequest,
					Message: "invalid parameters",
					Errors:  fieldsErrors,
				},
			})
			return
		}

		// Log before dispatch
		log.Printf("Transfer request from %s to %s for %d by user %s", from, to, amount, userID)

		// Dispatch
		err = act.MoneyTransfer(from, to, amount)
		if err != nil {
			if errors.Is(err, entities.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(Response{
					Error: Error{
						Code:    http.StatusNotFound,
						Message: err.Error(),
					},
				})
				return
			} else if e, ok := err.(*entities.ErrInvalidArgument); ok {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(Response{
					Error: Error{
						Code:    http.StatusBadRequest,
						Message: e.Error(),
						Errors: []ErrorItem{
							{
								Message:      e.Error(),
								Reason:       e.Reason(),
								LocationType: LocationTypeParameter,
								Location:     e.Argument,
							},
						},
					},
				})
				return
			} else if e, ok := err.(*entities.ErrInvalidState); ok {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(Response{
					Error: Error{
						Code:    http.StatusOK,
						Message: e.Error(),
						Errors: []ErrorItem{
							{
								Message:      e.Error(),
								Reason:       e.Reason(),
								LocationType: LocationTypeParameter,
								Location:     e.RelatedArgument,
							},
						},
					},
				})
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(Response{
					Error: Error{
						Code:    http.StatusInternalServerError,
						Message: err.Error(),
					},
				})
				return
			}
		}

		// Success
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			Data: nil,
		})
	})
}
