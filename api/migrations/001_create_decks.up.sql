CREATE TABLE IF NOT EXISTS decks (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    parent_id BIGINT REFERENCES decks(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    algorithm VARCHAR(50) NOT NULL DEFAULT 'sm2',
    srs_config JSONB NOT NULL DEFAULT '{
        "sm2": {
            "initial_ease_factor": 2.5,
            "min_ease_factor": 1.3,
            "max_ease_factor": 3.0,
            "ease_decrement": 0.2,
            "ease_increment": 0.15,
            "easy_bonus_multiplier": 1.3,
            "graduating_interval": 1,
            "mastered_threshold": 21
        }
    }'::jsonb,
    position INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, parent_id, name)
);

CREATE INDEX idx_decks_user_id ON decks(user_id);
CREATE INDEX idx_decks_parent_id ON decks(parent_id);
CREATE INDEX idx_decks_algorithm ON decks(algorithm);

COMMENT ON COLUMN decks.algorithm IS 'SRS algorithm: sm2, fsrs';
COMMENT ON COLUMN decks.srs_config IS 'Algorithm-specific configuration as JSON. For SM-2: {sm2: {initial_ease_factor, min_ease_factor, max_ease_factor, ease_decrement, ease_increment, easy_bonus_multiplier, graduating_interval, mastered_threshold}}';
