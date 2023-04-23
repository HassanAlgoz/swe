package http

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/internal/common"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (c *Server) registerHandlers() {
	// Prometheus Metrics
	c.mux.Handle("/metrics", promhttp.Handler())

	// Action handlers
	c.mux.HandleFunc("/transfer:transfer-money", c.TransferMoney)
	c.mux.HandleFunc("/user:get-account", c.GetAccount)
}

func (c *Server) TransferMoney(w http.ResponseWriter, r *http.Request) {
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
	headers := requireHeaders(w, r, []string{
		"x-user-id",
	})
	userID := headers["x-user-id"]
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

	// pre-action Log
	log.Printf("MoneyTransfer from %s to %s for %d by user %s", from, to, amount, userID)

	// Invoke the action
	err = c.app.MoneyTransfer(from, to, amount)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			ErrNotFound(w, err)
		} else if e, ok := err.(*common.ErrInvalidArgument); ok {
			ErrInvalidArgument(w, e)
		} else if e, ok := err.(*common.ErrInvalidState); ok {
			ErrInvalidState(w, e)
		} else {
			ErrInternal(w, err)
		}
		return
	}

	// Success
	Ok(w, nil)
}

func (c *Server) GetAccount(w http.ResponseWriter, r *http.Request) {
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
		ID string `json:"id"`
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
	id, err := uuid.Parse(fields.ID)
	if err != nil {
		fieldsErrors = append(fieldsErrors, ErrorItem{
			LocationType: "parameter",
			Location:     "from",
			Message:      "invalid uuid",
			Reason:       err.Error(),
		})
	}

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

	// pre-action Log
	log.Printf("GetAccount %s by user %s", id, userID)

	// Invoke the action
	result, err := c.app.GetAccount(id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			ErrNotFound(w, err)
		} else {
			ErrInternal(w, err)
		}
		return
	}

	// Success
	Ok(w, result)
}
