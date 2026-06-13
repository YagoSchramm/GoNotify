package dto

import (
	"time"

	"github.com/YagoSchramm/GoNotify/domain/entity"
	"github.com/google/uuid"
)

type CreateTriggerRequest struct {
	EventType       string                 `json:"event_type"`
	DestinationType entity.DestinationType `json:"destination_type"`
	Config          map[string]any         `json:"config"`
}

type UpdateTriggerRequest struct {
	EventType *string        `json:"event_type,omitempty"`
	Config    map[string]any `json:"config,omitempty"`
	Active    *bool          `json:"active,omitempty"`
}

type TriggerResponse struct {
	ID              uuid.UUID              `json:"id"`
	EventType       string                 `json:"event_type"`
	DestinationType entity.DestinationType `json:"destination_type"`
	Config          map[string]any         `json:"config"`
	Active          bool                   `json:"active"`
	CreatedAt       time.Time              `json:"created_at"`
}

func ToTriggerResponse(t *entity.Trigger) *TriggerResponse {
	return &TriggerResponse{
		ID:              t.ID,
		EventType:       t.EventType,
		DestinationType: t.DestinationType,
		Config:          t.Config,
		Active:          t.Active,
		CreatedAt:       t.CreatedAt,
	}
}
