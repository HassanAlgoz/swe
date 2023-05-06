package port

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	grpcInbound "github.com/hassanalgoz/swe/pkg/inbound/grpc"
	inbound "github.com/hassanalgoz/swe/pkg/inbound/http"
	lmsPort "github.com/hassanalgoz/swe/pkg/services/ports/lms"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

	// Headers
	authToken, err := inbound.GetAuthToken(r)
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

	reqID, ok := inbound.GetRequestId(r)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(inbound.Response{
			Error: inbound.Error{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("Missing request header: %s", inbound.HeaderRequestId),
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
	log.Debug().Msgf("MoneyTransfer from %s to %s for %d by user %s", from, to, amount, reqID)

	// Embed http headers as gRPC headers
	md := metadata.New(map[string]string{
		string(grpcInbound.HeaderUserId):    authToken,
		string(grpcInbound.HeaderRequestId): reqID,
	})
	grpc.SendHeader(s.ctx, md)

	// Invoke the action
	id, err := lmsClient.CreateCourse(s.ctx, &lmsPort.CreateCourseRequest{
		Name:        from.String(),
		Description: "asdf asdf",
		Instructors: []string{"vczxcv", "123asdf"},
	})
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			statusCode := HTTPStatusFromCode(statusErr.Code())
			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(inbound.Response{
				Data: statusErr.Message(),
			})
		} else {
			log.Error().Msgf("Error while calling CreateCourse RPC: %v", err)
			inbound.ErrInternal(w, err)
		}
		return
	}

	// Success
	inbound.Ok(w, map[string]any{
		"id": id.String(),
	})
}
