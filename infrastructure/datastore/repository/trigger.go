package repository

import (
	"context"

	"github.com/YagoSchramm/GoNotify/domain/entity"
	"github.com/google/uuid"
)

type TriggerRepository interface {
	Create(ctx context.Context, trigger *entity.Trigger) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Trigger, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Trigger, error)
	Update(ctx context.Context, trigger *entity.Trigger) error
	Delete(ctx context.Context, id uuid.UUID) error
}
