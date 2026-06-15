// infrastructure/datastore/repository/impl/delivery_attempt_repository.go
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

//go:embed _query/delivery_attempt/create_delivery_attempt.sql
var createDeliveryAttemptQuery string

//go:embed _query/delivery_attempt/find_delivery_attempts_by_notification_id.sql
var findDeliveryAttemptsByNotificationIDQuery string

//go:embed _query/delivery_attempt/count_delivery_attempts_by_notification_id.sql
var countDeliveryAttemptsByNotificationIDQuery string

type deliveryAttemptRepository struct {
	db *sql.DB
}

func NewDeliveryAttemptRepository(db *sql.DB) repository.DeliveryAttemptRepository {
	return &deliveryAttemptRepository{db: db}
}

func (r *deliveryAttemptRepository) Create(ctx context.Context, a *entity.DeliveryAttempt) error {
	_, err := r.db.ExecContext(
		ctx,
		createDeliveryAttemptQuery,
		a.ID,
		a.NotificationID,
		a.AttemptNumber,
		string(a.Outcome),
		a.Error,
		a.AttemptedAt,
	)
	if err != nil {
		return derr.JoinError("failed to create delivery attempt", err)
	}

	return nil
}

func (r *deliveryAttemptRepository) FindByNotificationID(ctx context.Context, notificationID uuid.UUID) ([]*entity.DeliveryAttempt, error) {
	rows, err := r.db.QueryContext(ctx, findDeliveryAttemptsByNotificationIDQuery, notificationID)
	if err != nil {
		return nil, derr.JoinError("failed to find delivery attempts by notification id", err)
	}
	defer rows.Close()

	var attempts []*entity.DeliveryAttempt

	for rows.Next() {
		var a entity.DeliveryAttempt
		var outcome string

		err := rows.Scan(
			&a.ID,
			&a.NotificationID,
			&a.AttemptNumber,
			&outcome,
			&a.Error,
			&a.AttemptedAt,
		)
		if err != nil {
			return nil, derr.JoinError("failed to scan delivery attempt", err)
		}

		a.Outcome = entity.Outcome(outcome)
		attempts = append(attempts, &a)
	}

	if err := rows.Err(); err != nil {
		return nil, derr.JoinError("failed to iterate delivery attempts", err)
	}

	return attempts, nil
}

func (r *deliveryAttemptRepository) CountByNotificationID(ctx context.Context, notificationID uuid.UUID) (int32, error) {
	row := r.db.QueryRowContext(ctx, countDeliveryAttemptsByNotificationIDQuery, notificationID)

	var count int32
	if err := row.Scan(&count); err != nil {
		return 0, derr.JoinError("failed to count delivery attempts", err)
	}

	return count, nil
}
