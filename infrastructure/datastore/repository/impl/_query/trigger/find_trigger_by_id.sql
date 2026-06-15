-- _query/trigger/find_trigger_by_id.sql
SELECT id, user_id, event_type, destination_type, config, active, created_at
FROM triggers
WHERE id = $1