-- _query/trigger/update_trigger.sql
UPDATE triggers
SET event_type = $2,
    config = $3,
    active = $4
WHERE id = $1