package impl

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/YagoSchramm/GoNotify/domain/entity"
	"github.com/YagoSchramm/GoNotify/domain/entity/derr"
	"github.com/YagoSchramm/GoNotify/infrastructure/datastore/repository"
	"github.com/google/uuid"
)

//go:embed _query/notification/create_notification.sql
var createNotificationQuery string

//go:embed _query/notification/find_notification_by_id.sql
var findNotificationByIDQuery string

//go:embed _query/notification/find_notifications_by_trigger_id.sql
var findNotificationsByTriggerIDQuery string

//go:embed _query/notification/find_notification_by_idempotency_key.sql
var findNotificationByIdempotencyKeyQuery string

//go:embed _query/notification/update_notification_status.sql
var updateNotificationStatusQuery string

type notificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) repository.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, n *entity.Notification) error {
	_, err := r.db.ExecContext(
		ctx,
		createNotificationQuery,
		n.ID,
		n.TriggerID,
		string(n.Status),
		n.IdempotencyKey,
		n.CreatedAt,
	)
	if err != nil {
		return derr.JoinError("failed to create notification", err)
	}

	return nil
}

func (r *notificationRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Notification, error) {
	row := r.db.QueryRowContext(ctx, findNotificationByIDQuery, id)

	var n entity.Notification
	var status string

	err := row.Scan(
		&n.ID,
		&n.TriggerID,
		&status,
		&n.IdempotencyKey,
		&n.CreatedAt,
		&n.DeliveredAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, derr.NotFoundError
		}
		return nil, derr.JoinError("failed to find notification by id", err)
	}

	n.Status = entity.Status(status)

	return &n, nil
}

func (r *notificationRepository) FindByTriggerID(ctx context.Context, triggerID uuid.UUID) ([]*entity.Notification, error) {
	rows, err := r.db.QueryContext(ctx, findNotificationsByTriggerIDQuery, triggerID)
	if err != nil {
		return nil, derr.JoinError("failed to find notifications by trigger id", err)
	}
	defer rows.Close()

	var notifications []*entity.Notification

	for rows.Next() {
		var n entity.Notification
		var status string

		err := rows.Scan(
			&n.ID,
			&n.TriggerID,
			&status,
			&n.IdempotencyKey,
			&n.CreatedAt,
			&n.DeliveredAt,
		)
		if err != nil {
			return nil, derr.JoinError("failed to scan notification", err)
		}

		n.Status = entity.Status(status)
		notifications = append(notifications, &n)
	}

	if err := rows.Err(); err != nil {
		return nil, derr.JoinError("failed to iterate notifications", err)
	}

	return notifications, nil
}

func (r *notificationRepository) FindByIdempotencyKey(ctx context.Context, key string) (*entity.Notification, error) {
	row := r.db.QueryRowContext(ctx, findNotificationByIdempotencyKeyQuery, key)

	var n entity.Notification
	var status string

	err := row.Scan(
		&n.ID,
		&n.TriggerID,
		&status,
		&n.IdempotencyKey,
		&n.CreatedAt,
		&n.DeliveredAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, derr.NotFoundError
		}
		return nil, derr.JoinError("failed to find notification by idempotency key", err)
	}

	n.Status = entity.Status(status)

	return &n, nil
}

func (r *notificationRepository) UpdateStatus(ctx context.Context, n *entity.Notification) error {
	_, err := r.db.ExecContext(
		ctx,
		updateNotificationStatusQuery,
		n.ID,
		string(n.Status),
		n.DeliveredAt,
	)
	if err != nil {
		return derr.JoinError("failed to update notification status", err)
	}

	return nil
}
