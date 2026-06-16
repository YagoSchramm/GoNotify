package rules

import (
	"errors"
	"fmt"

	"github.com/YagoSchramm/GoNotify/domain/entity"
)

func ValidateEventType(eventType string) error {
	if eventType == "" {
		return errors.New("event_type cannot be empty")
	}
	return nil
}
func ValidateConfig(t *entity.Trigger) error {
	switch t.DestinationType {
	case entity.DestinationEmail:
		to, ok := t.Config["to"].(string)
		if !ok || to == "" {
			return errors.New("email config requires 'to'")
		}
	case entity.DestinationWebhook:
		url, ok := t.Config["url"].(string)
		if !ok || url == "" {
			return errors.New("webhook config requires 'url'")
		}
	default:
		return fmt.Errorf("unknown destination type: %s", t.DestinationType)
	}
	return nil
}
