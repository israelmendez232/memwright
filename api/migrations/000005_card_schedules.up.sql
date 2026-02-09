CREATE TYPE schedule_state AS ENUM (
    'new',
    'learning',
    'review',
    'relearning',
    'mastered'
);

CREATE TABLE card_schedules (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    card_id UUID NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    state schedule_state NOT NULL DEFAULT 'new',
    ease_factor DECIMAL(4,2) NOT NULL DEFAULT 2.50,
    interval_days INT NOT NULL DEFAULT 0,
    repetitions INT NOT NULL DEFAULT 0,
    stability DECIMAL(10,4),
    difficulty DECIMAL(4,2),
    due_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_reviewed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, card_id)
);

CREATE INDEX idx_card_schedules_user_id ON card_schedules(user_id);
CREATE INDEX idx_card_schedules_card_id ON card_schedules(card_id);
CREATE INDEX idx_card_schedules_due_at ON card_schedules(due_at);

CREATE TRIGGER set_card_schedules_updated_at
    BEFORE UPDATE ON card_schedules
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();
