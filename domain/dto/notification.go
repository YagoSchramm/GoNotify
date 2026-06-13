package dto

import (
	"time"

	"github.com/YagoSchramm/GoNotify/domain/entity"
	"github.com/google/uuid"
)

type FireNotificationRequest struct {
	TriggerID      uuid.UUID `json:"trigger_id"`
	IdempotencyKey string    `json:"idempotency_key"`
}

type NotificationResponse struct {
	ID             uuid.UUID     `json:"id"`
	TriggerID      uuid.UUID     `json:"trigger_id"`
	Status         entity.Status `json:"status"`
	IdempotencyKey string        `json:"idempotency_key"`
	CreatedAt      time.Time     `json:"created_at"`
	DeliveredAt    *time.Time    `json:"delivered_at,omitempty"`
}

func ToNotificationResponse(n *entity.Notification) *NotificationResponse {
	return &NotificationResponse{
		ID:             n.ID,
		TriggerID:      n.TriggerID,
		Status:         n.Status,
		IdempotencyKey: n.IdempotencyKey,
		CreatedAt:      n.CreatedAt,
		DeliveredAt:    n.DeliveredAt,
	}
}
