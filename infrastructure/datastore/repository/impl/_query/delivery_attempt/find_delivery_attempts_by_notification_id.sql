SELECT id, notification_id, attempt_number, outcome, error, attempted_at
FROM delivery_attempts
WHERE notification_id = $1
ORDER BY attempt_number ASC