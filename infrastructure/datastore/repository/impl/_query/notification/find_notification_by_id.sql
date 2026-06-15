SELECT id, trigger_id, status, idempotency_key, created_at, delivered_at
FROM notifications
WHERE id = $1