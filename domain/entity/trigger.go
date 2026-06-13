package entity

import (
	"time"

	"github.com/google/uuid"
)

type DestinationType string

const (
	DestinationEmail   DestinationType = "email"
	DestinationWebhook DestinationType = "webhook"
)

type EmailConfig struct {
	To              string `json:"to"`
	SubjectTemplate string `json:"subject_template"`
	BodyTemplate    string `json:"body_template"`
}

type WebhookConfig struct {
	URL             string            `json:"url"`
	Method          string            `json:"method"`
	Headers         map[string]string `json:"headers,omitempty"`
	PayloadTemplate string            `json:"payload_template"`
}

type Trigger struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	EventType       string
	DestinationType DestinationType
	Config          map[string]any
	Active          bool
	CreatedAt       time.Time
}
