-- _query/trigger/create_trigger.sql
INSERT INTO triggers (id, user_id, event_type, destination_type, config, active, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)