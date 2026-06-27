package impl

import (
	"context"
	"errors"
	"testing"

	"github.com/YagoSchramm/GoNotify/domain/dto"
	"github.com/YagoSchramm/GoNotify/domain/entity"
	"github.com/YagoSchramm/GoNotify/domain/entity/derr"
	"github.com/google/uuid"
)

func TestNotificationUseCaseFireAndFindByTriggerID(t *testing.T) {
	triggerID := uuid.New()
	userID := uuid.New()
	triggerRepo := &fakeTriggerRepository{triggers: []*entity.Trigger{{ID: triggerID, UserID: userID, Active: true}}}
	notificationRepo := &fakeNotificationRepository{}
	uc := NewNotificationUseCase(notificationRepo, triggerRepo)

	resp, err := uc.Fire(context.Background(), userID, dto.FireNotificationRequest{TriggerID: triggerID, IdempotencyKey: "key-1"})
	if err != nil {
		t.Fatalf("Fire returned error: %v", err)
	}
	if resp == nil || resp.TriggerID != triggerID {
		t.Fatal("expected a notification response for the active trigger")
	}
	if len(notificationRepo.notifications) != 1 {
		t.Fatalf("expected one notification to be created, got %d", len(notificationRepo.notifications))
	}

	duplicate, err := uc.Fire(context.Background(), userID, dto.FireNotificationRequest{TriggerID: triggerID, IdempotencyKey: "key-1"})
	if err != nil {
		t.Fatalf("duplicate Fire returned error: %v", err)
	}
	if duplicate.ID != resp.ID {
		t.Fatal("expected duplicate request to return the existing notification")
	}
	if len(notificationRepo.notifications) != 1 {
		t.Fatalf("expected no duplicate notification to be stored, got %d", len(notificationRepo.notifications))
	}
}

func TestNotificationUseCaseFindByIDRequiresOwnerMatch(t *testing.T) {
	triggerID := uuid.New()
	ownerID := uuid.New()
	otherUserID := uuid.New()
	triggerRepo := &fakeTriggerRepository{triggers: []*entity.Trigger{{ID: triggerID, UserID: ownerID, Active: true}}}
	notificationRepo := &fakeNotificationRepository{}
	if err := notificationRepo.Create(context.Background(), &entity.Notification{ID: uuid.New(), TriggerID: triggerID, IdempotencyKey: "key-2"}); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	uc := NewNotificationUseCase(notificationRepo, triggerRepo)
	_, err := uc.FindByID(context.Background(), otherUserID, notificationRepo.notifications[0].ID)
	if !errors.Is(err, derr.NotFoundError) {
		t.Fatalf("expected not found for unrelated user, got %v", err)
	}
}
