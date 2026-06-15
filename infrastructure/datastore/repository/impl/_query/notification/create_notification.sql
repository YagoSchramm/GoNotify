INSERT INTO notifications (id, trigger_id, status, idempotency_key, created_at)
VALUES ($1, $2, $3, $4, $5)