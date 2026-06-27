package impl

import (
	"context"

	"github.com/YagoSchramm/GoNotify/domain/entity"
	"github.com/YagoSchramm/GoNotify/domain/entity/derr"
	"github.com/google/uuid"
)

type fakeAuthRepository struct {
	users []*entity.User
}

func (f *fakeAuthRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	for _, user := range f.users {
		if user.Email == email {
			return cloneUser(user), nil
		}
	}
	return nil, derr.NotFoundError
}

func (f *fakeAuthRepository) AttemptRegister(ctx context.Context, user entity.User) (*uuid.UUID, error) {
	id := uuid.New()
	f.users = append(f.users, &entity.User{ID: id, Email: user.Email, Password: user.Password})
	return &id, nil
}

func (f *fakeAuthRepository) AttemptLogin(ctx context.Context, credentials entity.UserCredentials) (*entity.User, error) {
	for _, user := range f.users {
		if user.Email == credentials.Email {
			return cloneUser(user), nil
		}
	}
	return nil, derr.NotFoundError
}

type fakeNotificationRepository struct {
	notifications []*entity.Notification
}

func (f *fakeNotificationRepository) Create(ctx context.Context, notification *entity.Notification) error {
	f.notifications = append(f.notifications, cloneNotification(notification))
	return nil
}

func (f *fakeNotificationRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Notification, error) {
	for _, notification := range f.notifications {
		if notification.ID == id {
			return cloneNotification(notification), nil
		}
	}
	return nil, derr.NotFoundError
}

func (f *fakeNotificationRepository) FindByTriggerID(ctx context.Context, triggerID uuid.UUID) ([]*entity.Notification, error) {
	var results []*entity.Notification
	for _, notification := range f.notifications {
		if notification.TriggerID == triggerID {
			results = append(results, cloneNotification(notification))
		}
	}
	return results, nil
}

func (f *fakeNotificationRepository) FindByIdempotencyKey(ctx context.Context, key string) (*entity.Notification, error) {
	for _, notification := range f.notifications {
		if notification.IdempotencyKey == key {
			return cloneNotification(notification), nil
		}
	}
	return nil, derr.NotFoundError
}

func (f *fakeNotificationRepository) UpdateStatus(ctx context.Context, notification *entity.Notification) error {
	for i, existing := range f.notifications {
		if existing.ID == notification.ID {
			f.notifications[i] = cloneNotification(notification)
			return nil
		}
	}
	return derr.NotFoundError
}

type fakeTriggerRepository struct {
	triggers []*entity.Trigger
}

func (f *fakeTriggerRepository) Create(ctx context.Context, trigger *entity.Trigger) error {
	f.triggers = append(f.triggers, cloneTrigger(trigger))
	return nil
}

func (f *fakeTriggerRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Trigger, error) {
	for _, trigger := range f.triggers {
		if trigger.ID == id {
			return cloneTrigger(trigger), nil
		}
	}
	return nil, derr.NotFoundError
}

func (f *fakeTriggerRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Trigger, error) {
	var results []*entity.Trigger
	for _, trigger := range f.triggers {
		if trigger.UserID == userID {
			results = append(results, cloneTrigger(trigger))
		}
	}
	return results, nil
}

func (f *fakeTriggerRepository) Update(ctx context.Context, trigger *entity.Trigger) error {
	for i, existing := range f.triggers {
		if existing.ID == trigger.ID {
			f.triggers[i] = cloneTrigger(trigger)
			return nil
		}
	}
	return derr.NotFoundError
}

func (f *fakeTriggerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	for i, trigger := range f.triggers {
		if trigger.ID == id {
			f.triggers = append(f.triggers[:i], f.triggers[i+1:]...)
			return nil
		}
	}
	return derr.NotFoundError
}

func cloneUser(user *entity.User) *entity.User {
	if user == nil {
		return nil
	}
	copy := *user
	return &copy
}

func cloneNotification(notification *entity.Notification) *entity.Notification {
	if notification == nil {
		return nil
	}
	copy := *notification
	return &copy
}

func cloneTrigger(trigger *entity.Trigger) *entity.Trigger {
	if trigger == nil {
		return nil
	}
	copy := *trigger
	return &copy
}
