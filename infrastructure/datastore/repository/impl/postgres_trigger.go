package impl

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"

	"github.com/YagoSchramm/GoNotify/domain/entity"
	"github.com/YagoSchramm/GoNotify/domain/entity/derr"
	"github.com/YagoSchramm/GoNotify/infrastructure/datastore/repository"
	"github.com/google/uuid"
)

func NewTriggerRepository(db *sql.DB) repository.TriggerRepository {
	return &triggerRepository{
		db: db,
	}
}

type triggerRepository struct {
	db *sql.DB
}

//go:embed _query/trigger/create_trigger.sql
var createTriggerQuery string

//go:embed _query/trigger/find_trigger_by_id.sql
var findTriggerByIDQuery string

//go:embed _query/trigger/find_triggers_by_user_id.sql
var findTriggersByUserIDQuery string

//go:embed _query/trigger/update_trigger.sql
var updateTriggerQuery string

//go:embed _query/trigger/delete_trigger.sql
var deleteTriggerQuery string

func (r *triggerRepository) Create(ctx context.Context, t *entity.Trigger) error {
	configJSON, err := json.Marshal(t.Config)
	if err != nil {
		return derr.JoinError("failed to marshal config", err)
	}

	_, err = r.db.ExecContext(
		ctx,
		createTriggerQuery,
		t.ID,
		t.UserID,
		t.EventType,
		string(t.DestinationType),
		configJSON,
		t.Active,
		t.CreatedAt,
	)
	if err != nil {
		return derr.JoinError("failed to create trigger", err)
	}

	return nil
}

func (r *triggerRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Trigger, error) {
	row := r.db.QueryRowContext(ctx, findTriggerByIDQuery, id)

	var t entity.Trigger
	var configJSON []byte
	var destinationType string

	err := row.Scan(
		&t.ID,
		&t.UserID,
		&t.EventType,
		&destinationType,
		&configJSON,
		&t.Active,
		&t.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, derr.NotFoundError
		}
		return nil, derr.JoinError("failed to find trigger by id", err)
	}

	if err := json.Unmarshal(configJSON, &t.Config); err != nil {
		return nil, derr.JoinError("failed to unmarshal config", err)
	}

	t.DestinationType = entity.DestinationType(destinationType)

	return &t, nil
}

func (r *triggerRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Trigger, error) {
	rows, err := r.db.QueryContext(ctx, findTriggersByUserIDQuery, userID)
	if err != nil {
		return nil, derr.JoinError("failed to find triggers by user id", err)
	}
	defer rows.Close()

	var triggers []*entity.Trigger

	for rows.Next() {
		var t entity.Trigger
		var configJSON []byte
		var destinationType string

		err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.EventType,
			&destinationType,
			&configJSON,
			&t.Active,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, derr.JoinError("failed to scan trigger", err)
		}

		if err := json.Unmarshal(configJSON, &t.Config); err != nil {
			return nil, derr.JoinError("failed to unmarshal config", err)
		}

		t.DestinationType = entity.DestinationType(destinationType)
		triggers = append(triggers, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, derr.JoinError("failed to iterate triggers", err)
	}

	return triggers, nil
}

func (r *triggerRepository) Update(ctx context.Context, t *entity.Trigger) error {
	configJSON, err := json.Marshal(t.Config)
	if err != nil {
		return derr.JoinError("failed to marshal config", err)
	}

	_, err = r.db.ExecContext(
		ctx,
		updateTriggerQuery,
		t.ID,
		t.EventType,
		configJSON,
		t.Active,
	)
	if err != nil {
		return derr.JoinError("failed to update trigger", err)
	}

	return nil
}

func (r *triggerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, deleteTriggerQuery, id)
	if err != nil {
		return derr.JoinError("failed to delete trigger", err)
	}

	return nil
}
