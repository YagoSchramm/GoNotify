package impl

import (
	"context"
	"errors"
	"time"

	"github.com/YagoSchramm/GoNotify/domain/dto"
	"github.com/YagoSchramm/GoNotify/domain/entity"
	"github.com/YagoSchramm/GoNotify/domain/entity/derr"
	"github.com/YagoSchramm/GoNotify/domain/rules"
	"github.com/YagoSchramm/GoNotify/domain/usecase"
	"github.com/YagoSchramm/GoNotify/infrastructure/datastore/repository"
	"github.com/google/uuid"
)

type triggerUseCase struct {
	repository repository.TriggerRepository
}

func NewTriggerUseCase(repository repository.TriggerRepository) usecase.TriggerUseCase {
	return &triggerUseCase{repository: repository}
}

func (u *triggerUseCase) Create(ctx context.Context, userID uuid.UUID, req dto.CreateTriggerRequest) (*dto.TriggerResponse, error) {
	t := &entity.Trigger{
		ID:              uuid.New(),
		UserID:          userID,
		EventType:       req.EventType,
		DestinationType: req.DestinationType,
		Config:          req.Config,
		Active:          true,
		CreatedAt:       time.Now(),
	}

	if err := rules.ValidateEventType(t.EventType); err != nil {
		return nil, err
	}

	if err := rules.ValidateConfig(t); err != nil {
		return nil, err
	}

	if err := u.repository.Create(ctx, t); err != nil {
		return nil, derr.JoinError("failed to create trigger", err)
	}

	return dto.ToTriggerResponse(t), nil
}

func (u *triggerUseCase) FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*dto.TriggerResponse, error) {
	t, err := u.repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, derr.NotFoundError) {
			return nil, derr.NotFoundError
		}
		return nil, derr.JoinError("failed to find trigger", err)
	}

	if t.UserID != userID {
		return nil, derr.NotFoundError
	}

	return dto.ToTriggerResponse(t), nil
}

func (u *triggerUseCase) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.TriggerResponse, error) {
	triggers, err := u.repository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, derr.JoinError("failed to find triggers", err)
	}

	var responses []*dto.TriggerResponse
	for _, t := range triggers {
		responses = append(responses, dto.ToTriggerResponse(t))
	}

	return responses, nil
}

func (u *triggerUseCase) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, req dto.UpdateTriggerRequest) (*dto.TriggerResponse, error) {
	t, err := u.repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, derr.NotFoundError) {
			return nil, derr.NotFoundError
		}
		return nil, derr.JoinError("failed to find trigger", err)
	}

	if t.UserID != userID {
		return nil, derr.NotFoundError
	}

	if req.EventType != nil {
		t.EventType = *req.EventType
	}
	if req.Config != nil {
		t.Config = req.Config
	}
	if req.Active != nil {
		t.Active = *req.Active
	}

	if err := rules.ValidateEventType(t.EventType); err != nil {
		return nil, err
	}

	if err := rules.ValidateConfig(t); err != nil {
		return nil, err
	}

	if err := u.repository.Update(ctx, t); err != nil {
		return nil, derr.JoinError("failed to update trigger", err)
	}

	return dto.ToTriggerResponse(t), nil
}

func (u *triggerUseCase) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	t, err := u.repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, derr.NotFoundError) {
			return derr.NotFoundError
		}
		return derr.JoinError("failed to find trigger", err)
	}

	if t.UserID != userID {
		return derr.NotFoundError
	}

	if err := u.repository.Delete(ctx, id); err != nil {
		return derr.JoinError("failed to delete trigger", err)
	}

	return nil
}
