package impl

import (
	"context"
	"errors"
	"testing"

	"github.com/YagoSchramm/GoNotify/domain/entity"
	"github.com/YagoSchramm/GoNotify/domain/entity/derr"
	"github.com/YagoSchramm/GoNotify/infrastructure/foundation/hash"
	"github.com/google/uuid"
)

func TestAuthUseCaseAttemptRegister(t *testing.T) {
	repo := &fakeAuthRepository{}
	uc := NewAuthRepository(repo, "super-secret")

	user := entity.User{Email: "new@example.com", Password: "Password123"}
	token, err := uc.AttemptRegister(context.Background(), user)
	if err != nil {
		t.Fatalf("AttemptRegister returned error: %v", err)
	}
	if token == "" {
		t.Fatal("expected a token to be generated")
	}
	if len(repo.users) != 1 {
		t.Fatalf("expected one stored user, got %d", len(repo.users))
	}
	if repo.users[0].Password == user.Password {
		t.Fatal("expected password to be hashed before storage")
	}
}

func TestAuthUseCaseAttemptLogin(t *testing.T) {
	repo := &fakeAuthRepository{}
	hashedPassword, err := hash.Hash("Password123")
	if err != nil {
		t.Fatalf("hash.Hash returned error: %v", err)
	}
	repo.users = append(repo.users, &entity.User{Email: "user@example.com", Password: hashedPassword})

	uc := NewAuthRepository(repo, "super-secret")
	credentials := entity.UserCredentials{Email: "user@example.com", Password: "Password123"}

	token, err := uc.AttemptLogin(context.Background(), credentials)
	if err != nil {
		t.Fatalf("AttemptLogin returned error: %v", err)
	}
	if token == "" {
		t.Fatal("expected a token to be generated")
	}
}

func TestAuthUseCaseAttemptLoginWithWrongPassword(t *testing.T) {
	repo := &fakeAuthRepository{}
	hashedPassword, err := hash.Hash("Password123")
	if err != nil {
		t.Fatalf("hash.Hash returned error: %v", err)
	}
	repo.users = append(repo.users, &entity.User{Email: "user@example.com", Password: hashedPassword})

	uc := NewAuthRepository(repo, "super-secret")
	credentials := entity.UserCredentials{Email: "user@example.com", Password: "WrongPassword1"}

	_, err = uc.AttemptLogin(context.Background(), credentials)
	if !errors.Is(err, derr.InvalidCredentials) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}

func TestAuthUseCaseValidateSession(t *testing.T) {
	repo := &fakeAuthRepository{}
	userID := uuid.New()
	repo.users = append(repo.users, &entity.User{ID: userID, Email: "session@example.com"})
	uc := NewAuthRepository(repo, "super-secret")

	if err := uc.ValidateSession(context.Background(), userID, "session@example.com"); err != nil {
		t.Fatalf("expected session validation to succeed, got %v", err)
	}

	if err := uc.ValidateSession(context.Background(), uuid.New(), "session@example.com"); !errors.Is(err, derr.UnauthorizedError) {
		t.Fatalf("expected unauthorized error, got %v", err)
	}
}
