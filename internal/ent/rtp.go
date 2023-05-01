package ent

import (
	"time"

	"github.com/google/uuid"
)

type RTPState string

const (
	RTPStateInitiated RTPState = "INITIATED"
	RTPStateRejected  RTPState = "REJECTED"
	RTPStateAccepted  RTPState = "ACCEPTED"
	RTPStateCancelled RTPState = "CANCELLED"
	RTPStateExpired   RTPState = "EXPIRED"
)

type RTP struct {
	ID        uuid.UUID
	Payee     uuid.UUID
	Payer     uuid.UUID
	Amount    uint
	CreatedAt time.Time
	Due       time.Time
	State     RTPState
}
