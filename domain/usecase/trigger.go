package usecase

import (
	"context"

	"github.com/YagoSchramm/GoNotify/domain/dto"
	"github.com/google/uuid"
)

type TriggerUseCase interface {
	Create(ctx context.Context, userID uuid.UUID, req dto.CreateTriggerRequest) (*dto.TriggerResponse, error)
	FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*dto.TriggerResponse, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.TriggerResponse, error)
	Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, req dto.UpdateTriggerRequest) (*dto.TriggerResponse, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
}
