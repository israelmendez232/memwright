CREATE TYPE review_rating AS ENUM (
    'wrong',
    'correct',
    'easy'
);

CREATE TABLE review_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    card_id UUID NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    rating review_rating NOT NULL,
    state_before schedule_state NOT NULL,
    state_after schedule_state NOT NULL,
    ease_factor_before DECIMAL(4,2) NOT NULL,
    ease_factor_after DECIMAL(4,2) NOT NULL,
    interval_before INT NOT NULL,
    interval_after INT NOT NULL,
    time_spent_ms INT,
    reviewed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_review_logs_user_id ON review_logs(user_id);
CREATE INDEX idx_review_logs_card_id ON review_logs(card_id);
CREATE INDEX idx_review_logs_reviewed_at ON review_logs(reviewed_at);
CREATE INDEX idx_review_logs_user_reviewed_at ON review_logs(user_id, reviewed_at);
