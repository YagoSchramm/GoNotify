package repository

import (
	"context"

	"github.com/YagoSchramm/GoNotify/domain/entity"
	"github.com/google/uuid"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *entity.Notification) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Notification, error)
	FindByTriggerID(ctx context.Context, triggerID uuid.UUID) ([]*entity.Notification, error)
	FindByIdempotencyKey(ctx context.Context, key string) (*entity.Notification, error)
	UpdateStatus(ctx context.Context, notification *entity.Notification) error
}
