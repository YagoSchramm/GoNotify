package rules

import (
	"errors"

	"github.com/YagoSchramm/GoNotify/domain/entity"
)

func ValidateIdempotencyKey(n *entity.Notification) error {
	if n.IdempotencyKey == "" {
		return errors.New("idempotency_key cannot be empty")
	}
	return nil
}

func CanFire(n *entity.Notification) error {
	if n.Status != entity.StatusPending {
		return errors.New("notification already processed")
	}
	return nil
}

func CanMarkDelivered(n *entity.Notification) error {
	if n.Status == entity.StatusDelivered {
		return errors.New("notification already delivered")
	}
	if n.Status == entity.StatusFailed {
		return errors.New("cannot deliver a failed notification")
	}
	return nil
}

func CanMarkFailed(n *entity.Notification) error {
	if n.Status == entity.StatusDelivered {
		return errors.New("cannot fail a delivered notification")
	}
	return nil
}
