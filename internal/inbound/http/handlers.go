package http

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/internal/ent"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *Server) registerHandlers() {
	// Prometheus Metrics
	s.mux.Handle("/metrics", promhttp.Handler())

	// Action handlers
	s.registerEndpoint([]string{http.MethodPost}, "/transfer:transfer-money", s.TransferMoney, &middlewareOptions{
		RequiredFeatureFlags: []FeatureFlag{FeatureFlagMoneyTransfer},
		RequiredHeaders:      []Header{HeaderAuthorization, HeaderRequestId},
	})
	s.registerEndpoint([]string{http.MethodGet}, "/user:get-account", s.GetAccount, &middlewareOptions{
		RequiredHeaders: []Header{HeaderAuthorization},
	})
}

func (s *Server) TransferMoney(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Header
	userID := getUserId(H(r, HeaderAuthorization))

	// Body
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
	err = s.app.MoneyTransfer(from, to, amount)
	if err != nil {
		if errors.Is(err, ent.ErrNotFound) {
			ErrNotFound(w, err)
		} else if e, ok := err.(*ent.ErrInvalidArgument); ok {
			ErrInvalidArgument(w, e)
		} else if e, ok := err.(*ent.ErrInvalidState); ok {
			ErrInvalidState(w, e)
		} else {
			ErrInternal(w, err)
		}
		return
	}

	// Success
	Ok(w, nil)
}

func (s *Server) GetAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Header
	userID := getUserId(r.Header.Get(string(HeaderAuthorization)))

	// Body
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
	result, err := s.app.GetAccount(id)
	if err != nil {
		if errors.Is(err, ent.ErrNotFound) {
			ErrNotFound(w, err)
		} else {
			ErrInternal(w, err)
		}
		return
	}

	// Success
	Ok(w, result)
}
