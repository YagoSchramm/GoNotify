package entity

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusDelivered Status = "delivered"
	StatusFailed    Status = "failed"
)

type Notification struct {
	ID             uuid.UUID
	TriggerID      uuid.UUID
	Status         Status
	IdempotencyKey string
	CreatedAt      time.Time
	DeliveredAt    *time.Time
}
