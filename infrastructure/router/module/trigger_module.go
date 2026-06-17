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

func NewTriggerModule(triggerUseCase usecase.TriggerUseCase) router.Module {
	return &triggerModule{
		triggerUseCase: triggerUseCase,
		name:           "Trigger",
		path:           "/triggers",
	}
}

type triggerModule struct {
	triggerUseCase usecase.TriggerUseCase
	name           string
	path           string
}

func (m triggerModule) Name() string {
	return m.name
}

func (m triggerModule) Path() string {
	return m.path
}

func (m triggerModule) Middlewares() []mux.MiddlewareFunc {
	return nil
}

func (m triggerModule) Routes() []router.RouteDefinition {
	return []router.RouteDefinition{
		{
			Path:        "",
			Description: "Create a trigger",
			Handler:     m.create,
			HttpMethods: []string{http.MethodPost},
			Public:      false,
		},
		{
			Path:        "",
			Description: "List triggers for the authenticated user",
			Handler:     m.list,
			HttpMethods: []string{http.MethodGet},
			Public:      false,
		},
		{
			Path:        "/{id}",
			Description: "Get a trigger by id",
			Handler:     m.getByID,
			HttpMethods: []string{http.MethodGet},
			Public:      false,
		},
		{
			Path:        "/{id}",
			Description: "Update a trigger",
			Handler:     m.update,
			HttpMethods: []string{http.MethodPatch},
			Public:      false,
		},
		{
			Path:        "/{id}",
			Description: "Delete a trigger",
			Handler:     m.delete,
			HttpMethods: []string{http.MethodDelete},
			Public:      false,
		},
	}
}

func (m triggerModule) create(w http.ResponseWriter, r *http.Request) {
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

	var req dto.CreateTriggerRequest
	if err := json.Unmarshal(body, &req); err != nil {
		slog.ErrorContext(ctx, "failed to unmarshal request body", "error", err)
		router.HandleError(w, err)
		return
	}

	trigger, err := m.triggerUseCase.Create(ctx, userID, req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create trigger", "error", err)
		router.HandleError(w, err)
		return
	}

	if err := router.Write(w, trigger); err != nil {
		slog.ErrorContext(ctx, "failed to write response", "error", err)
	}
}

func (m triggerModule) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := util.UserIDFromContext(r)
	if err != nil {
		router.HandleError(w, err)
		return
	}

	triggers, err := m.triggerUseCase.FindByUserID(ctx, userID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to list triggers", "error", err)
		router.HandleError(w, err)
		return
	}

	if err := router.Write(w, triggers); err != nil {
		slog.ErrorContext(ctx, "failed to write response", "error", err)
	}
}

func (m triggerModule) getByID(w http.ResponseWriter, r *http.Request) {
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

	trigger, err := m.triggerUseCase.FindByID(ctx, userID, id)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get trigger", "error", err)
		router.HandleError(w, err)
		return
	}

	if err := router.Write(w, trigger); err != nil {
		slog.ErrorContext(ctx, "failed to write response", "error", err)
	}
}

func (m triggerModule) update(w http.ResponseWriter, r *http.Request) {
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read request body", "error", err)
		router.HandleError(w, err)
		return
	}

	var req dto.UpdateTriggerRequest
	if err := json.Unmarshal(body, &req); err != nil {
		slog.ErrorContext(ctx, "failed to unmarshal request body", "error", err)
		router.HandleError(w, err)
		return
	}

	trigger, err := m.triggerUseCase.Update(ctx, userID, id, req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to update trigger", "error", err)
		router.HandleError(w, err)
		return
	}

	if err := router.Write(w, trigger); err != nil {
		slog.ErrorContext(ctx, "failed to write response", "error", err)
	}
}

func (m triggerModule) delete(w http.ResponseWriter, r *http.Request) {
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

	if err := m.triggerUseCase.Delete(ctx, userID, id); err != nil {
		slog.ErrorContext(ctx, "failed to delete trigger", "error", err)
		router.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
