package impl

import (
	"context"
	"errors"
	"time"

	"github.com/YagoSchramm/GoNotify/domain/dto"
	"github.com/YagoSchramm/GoNotify/domain/entity"
	"github.com/YagoSchramm/GoNotify/domain/entity/derr"
	"github.com/YagoSchramm/GoNotify/domain/usecase"
	"github.com/YagoSchramm/GoNotify/infrastructure/datastore/repository"
	"github.com/google/uuid"
)

type notificationUseCase struct {
	notificationRepo repository.NotificationRepository
	triggerRepo      repository.TriggerRepository
}

func NewNotificationUseCase(
	notificationRepo repository.NotificationRepository,
	triggerRepo repository.TriggerRepository,
) usecase.NotificationUseCase {
	return &notificationUseCase{
		notificationRepo: notificationRepo,
		triggerRepo:      triggerRepo,
	}
}

func (u *notificationUseCase) Fire(ctx context.Context, userID uuid.UUID, req dto.FireNotificationRequest) (*dto.NotificationResponse, error) {
	existing, err := u.notificationRepo.FindByIdempotencyKey(ctx, req.IdempotencyKey)
	if err != nil && !errors.Is(err, derr.NotFoundError) {
		return nil, derr.JoinError("failed to check idempotency key", err)
	}
	if existing != nil {
		return dto.ToNotificationResponse(existing), nil
	}

	trigger, err := u.triggerRepo.FindByID(ctx, req.TriggerID)
	if err != nil {
		if errors.Is(err, derr.NotFoundError) {
			return nil, derr.NotFoundError
		}
		return nil, derr.JoinError("failed to find trigger", err)
	}

	if trigger.UserID != userID {
		return nil, derr.NotFoundError
	}

	if !trigger.Active {
		return nil, derr.NewBadRequestError("trigger is not active")
	}

	n := &entity.Notification{
		ID:             uuid.New(),
		TriggerID:      req.TriggerID,
		Status:         entity.StatusPending,
		IdempotencyKey: req.IdempotencyKey,
		CreatedAt:      time.Now(),
	}

	if err := u.notificationRepo.Create(ctx, n); err != nil {
		return nil, derr.JoinError("failed to create notification", err)
	}

	return dto.ToNotificationResponse(n), nil
}

func (u *notificationUseCase) FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*dto.NotificationResponse, error) {
	n, err := u.notificationRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, derr.NotFoundError) {
			return nil, derr.NotFoundError
		}
		return nil, derr.JoinError("failed to find notification", err)
	}

	trigger, err := u.triggerRepo.FindByID(ctx, n.TriggerID)
	if err != nil {
		if errors.Is(err, derr.NotFoundError) {
			return nil, derr.NotFoundError
		}
		return nil, derr.JoinError("failed to find trigger", err)
	}

	if trigger.UserID != userID {
		return nil, derr.NotFoundError
	}

	return dto.ToNotificationResponse(n), nil
}

func (u *notificationUseCase) FindByTriggerID(ctx context.Context, userID uuid.UUID, triggerID uuid.UUID) ([]*dto.NotificationResponse, error) {
	trigger, err := u.triggerRepo.FindByID(ctx, triggerID)
	if err != nil {
		if errors.Is(err, derr.NotFoundError) {
			return nil, derr.NotFoundError
		}
		return nil, derr.JoinError("failed to find trigger", err)
	}

	if trigger.UserID != userID {
		return nil, derr.NotFoundError
	}

	notifications, err := u.notificationRepo.FindByTriggerID(ctx, triggerID)
	if err != nil {
		return nil, derr.JoinError("failed to find notifications", err)
	}

	var responses []*dto.NotificationResponse
	for _, n := range notifications {
		responses = append(responses, dto.ToNotificationResponse(n))
	}

	return responses, nil
}
