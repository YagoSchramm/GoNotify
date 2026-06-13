package entity

import (
	"time"

	"github.com/google/uuid"
)

type Outcome string

const (
	OutcomeFailure Outcome = "failure"
	OutcomeSuccess Outcome = "success"
)

type DeliveryAttempt struct {
	ID             uuid.UUID
	NotificationID uuid.UUID
	AttemptNumber  int32
	Outcome        Outcome
	Error          *string
	AttemptedAt    time.Time
}
