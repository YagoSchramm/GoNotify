package repository

import (
	"context"

	"github.com/YagoSchramm/GoNotify/domain/entity"
	"github.com/google/uuid"
)

type DeliveryAttemptRepository interface {
	Create(ctx context.Context, attempt *entity.DeliveryAttempt) error
	FindByNotificationID(ctx context.Context, notificationID uuid.UUID) ([]*entity.DeliveryAttempt, error)
	CountByNotificationID(ctx context.Context, notificationID uuid.UUID) (int32, error)
}
