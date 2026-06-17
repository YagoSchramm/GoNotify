package module

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/YagoSchramm/GoNotify/domain/dto"
	"github.com/YagoSchramm/GoNotify/domain/usecase"
	"github.com/YagoSchramm/GoNotify/infrastructure/router"
	"github.com/YagoSchramm/GoNotify/infrastructure/router/module/util"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func NewNotificationModule(notificationUseCase usecase.NotificationUseCase) router.Module {
	return &notificationModule{
		notificationUseCase: notificationUseCase,
		name:                "Notification",
		path:                "/notifications",
	}
}

type notificationModule struct {
	notificationUseCase usecase.NotificationUseCase
	name                string
	path                string
}

func (m notificationModule) Name() string {
	return m.name
}

func (m notificationModule) Path() string {
	return m.path
}

func (m notificationModule) Middlewares() []mux.MiddlewareFunc {
	return nil
}

func (m notificationModule) Routes() []router.RouteDefinition {
	return []router.RouteDefinition{
		{
			Path:        "",
			Description: "Fire a notification",
			Handler:     m.fire,
			HttpMethods: []string{http.MethodPost},
			Public:      false,
		},
		{
			Path:        "/{id}",
			Description: "Get a notification by id",
			Handler:     m.getByID,
			HttpMethods: []string{http.MethodGet},
			Public:      false,
		},
	}
}

func (m notificationModule) fire(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := util.UserIDFromContext(r)
	if err != nil {
		router.HandleError(w, err)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read request body", "error", err)
		router.HandleError(w, err)
		return
	}

	var req dto.FireNotificationRequest
	if err := json.Unmarshal(body, &req); err != nil {
		slog.ErrorContext(ctx, "failed to unmarshal request body", "error", err)
		router.HandleError(w, err)
		return
	}

	notification, err := m.notificationUseCase.Fire(ctx, userID, req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to fire notification", "error", err)
		router.HandleError(w, err)
		return
	}

	if err := router.Write(w, notification); err != nil {
		slog.ErrorContext(ctx, "failed to write response", "error", err)
	}
}

func (m notificationModule) getByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := util.UserIDFromContext(r)
	if err != nil {
		router.HandleError(w, err)
		return
	}

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		router.HandleError(w, err)
		return
	}

	notification, err := m.notificationUseCase.FindByID(ctx, userID, id)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get notification", "error", err)
		router.HandleError(w, err)
		return
	}

	if err := router.Write(w, notification); err != nil {
		slog.ErrorContext(ctx, "failed to write response", "error", err)
	}
}
