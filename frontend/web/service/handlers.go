package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/pkg/entities"
	inbound "github.com/hassanalgoz/swe/pkg/inbound/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *service) registerHandlers() {
	// Prometheus Metrics
	s.mux.Handle("/metrics", promhttp.Handler())

	// Action handlers
	s.registerEndpoint([]string{http.MethodPost}, "/transfer:transfer-money", s.TransferMoney, &endpointOptions{
		RequiredFeatureFlags: []string{"money_transfer"},
		RequiredHeaders:      []inbound.Header{inbound.HeaderAuthorization, inbound.HeaderRequestId},
	})
	// s.registerEndpoint([]string{http.MethodGet}, "/user:get-account", s.GetAccount, &inbound.MiddlewareOptions{
	// 	RequiredHeaders: []inbound.Header{inbound.HeaderAuthorization},
	// })
}

func (s *service) TransferMoney(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Header
	reqID, err := inbound.ExtractRequestId(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(inbound.Response{
			Error: inbound.Error{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("Invalid request header: %v", err),
			},
		})
		return
	}

	// Body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(inbound.Response{
			Error: inbound.Error{
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
		json.NewEncoder(w).Encode(inbound.Response{
			Error: inbound.Error{
				Code:    http.StatusBadRequest,
				Message: "Invalid json",
			},
		})
		return
	}

	// Parse body.json.fields
	// note: errors are appended
	var fieldsErrors []inbound.ErrorItem
	from, err := uuid.Parse(fields.From)
	if err != nil {
		fieldsErrors = append(fieldsErrors, inbound.ErrorItem{
			LocationType: "parameter",
			Location:     "from",
			Message:      "invalid uuid",
			Reason:       err.Error(),
		})
	}
	to, err := uuid.Parse(fields.To)
	if err != nil {
		fieldsErrors = append(fieldsErrors, inbound.ErrorItem{
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
		json.NewEncoder(w).Encode(inbound.Response{
			Error: inbound.Error{
				Code:    http.StatusBadRequest,
				Message: "invalid parameters",
				Errors:  fieldsErrors,
			},
		})
		return
	}

	// pre-action Log
	log.Printf("MoneyTransfer from %s to %s for %d by user %s", from, to, amount, reqID)

	// Invoke the action
	err = s.app.MoneyTransfer(from, to, amount)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			inbound.ErrNotFound(w, err)
		} else if e, ok := err.(*entities.ErrInvalidArgument); ok {
			inbound.ErrInvalidArgument(w, e)
		} else if e, ok := err.(*entities.ErrInvalidState); ok {
			inbound.ErrInvalidState(w, e)
		} else {
			inbound.ErrInternal(w, err)
		}
		return
	}

	// Success
	inbound.Ok(w, nil)
}
