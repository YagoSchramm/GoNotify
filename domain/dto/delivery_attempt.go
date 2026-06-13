package dto

import (
	"time"

	"github.com/YagoSchramm/GoNotify/domain/entity"
	"github.com/google/uuid"
)

type DeliveryAttemptResponse struct {
	ID             uuid.UUID      `json:"id"`
	NotificationID uuid.UUID      `json:"notification_id"`
	AttemptNumber  int32          `json:"attempt_number"`
	Outcome        entity.Outcome `json:"outcome"`
	Error          *string        `json:"error,omitempty"`
	AttemptedAt    time.Time      `json:"attempted_at"`
}

func ToDeliveryAttemptResponse(d *entity.DeliveryAttempt) *DeliveryAttemptResponse {
	return &DeliveryAttemptResponse{
		ID:             d.ID,
		NotificationID: d.NotificationID,
		AttemptNumber:  d.AttemptNumber,
		Outcome:        d.Outcome,
		Error:          d.Error,
		AttemptedAt:    d.AttemptedAt,
	}
}
