SELECT id, trigger_id, status, idempotency_key, created_at, delivered_at
FROM notifications
WHERE trigger_id = $1
ORDER BY created_at DESC