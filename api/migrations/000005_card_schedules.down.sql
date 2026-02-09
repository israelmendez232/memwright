DROP TRIGGER IF EXISTS set_card_schedules_updated_at ON card_schedules;
DROP INDEX IF EXISTS idx_card_schedules_due_at;
DROP INDEX IF EXISTS idx_card_schedules_card_id;
DROP INDEX IF EXISTS idx_card_schedules_user_id;
DROP TABLE IF EXISTS card_schedules;
DROP TYPE IF EXISTS schedule_state;
