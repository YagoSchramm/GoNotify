package impl

import (
	"context"
	"testing"

	"github.com/YagoSchramm/GoNotify/domain/dto"
	"github.com/YagoSchramm/GoNotify/domain/entity"
	"github.com/google/uuid"
)

func TestTriggerUseCaseCreateAndUpdate(t *testing.T) {
	triggerRepo := &fakeTriggerRepository{}
	uc := NewTriggerUseCase(triggerRepo)
	userID := uuid.New()

	created, err := uc.Create(context.Background(), userID, dto.CreateTriggerRequest{
		EventType:       "user.created",
		DestinationType: entity.DestinationEmail,
		Config:          map[string]any{"to": "ops@example.com"},
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if created == nil || created.Active != true {
		t.Fatal("expected the trigger to be created as active")
	}
	if len(triggerRepo.triggers) != 1 {
		t.Fatalf("expected one stored trigger, got %d", len(triggerRepo.triggers))
	}

	newEventType := "user.updated"
	newActive := false
	updated, err := uc.Update(context.Background(), userID, created.ID, dto.UpdateTriggerRequest{
		EventType: &newEventType,
		Config:    map[string]any{"to": "alerts@example.com"},
		Active:    &newActive,
	})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if updated.EventType != newEventType {
		t.Fatalf("expected event type to be updated, got %s", updated.EventType)
	}
	if updated.Active {
		t.Fatal("expected the trigger to be deactivated")
	}
}

func TestTriggerUseCaseDelete(t *testing.T) {
	triggerRepo := &fakeTriggerRepository{}
	uc := NewTriggerUseCase(triggerRepo)
	userID := uuid.New()
	trigger := &entity.Trigger{ID: uuid.New(), UserID: userID, EventType: "user.created", DestinationType: entity.DestinationEmail, Config: map[string]any{"to": "ops@example.com"}, Active: true}
	if err := triggerRepo.Create(context.Background(), trigger); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if err := uc.Delete(context.Background(), userID, trigger.ID); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if len(triggerRepo.triggers) != 0 {
		t.Fatalf("expected trigger to be removed, got %d", len(triggerRepo.triggers))
	}
}
