package usecase

import (
	"context"

	"github.com/YagoSchramm/GoNotify/domain/dto"
	"github.com/google/uuid"
)

type NotificationUseCase interface {
	Fire(ctx context.Context, userID uuid.UUID, req dto.FireNotificationRequest) (*dto.NotificationResponse, error)
	FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*dto.NotificationResponse, error)
	FindByTriggerID(ctx context.Context, userID uuid.UUID, triggerID uuid.UUID) ([]*dto.NotificationResponse, error)
}
