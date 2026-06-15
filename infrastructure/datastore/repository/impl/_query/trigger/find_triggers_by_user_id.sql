-- _query/trigger/find_triggers_by_user_id.sql
SELECT id, user_id, event_type, destination_type, config, active, created_at
FROM triggers
WHERE user_id = $1
ORDER BY created_at DESC