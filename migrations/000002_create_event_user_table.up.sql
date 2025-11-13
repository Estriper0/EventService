CREATE TABLE IF NOT EXISTS event.event_user (
    event_id INTEGER NOT NULL REFERENCES event.events(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    UNIQUE(event_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_event_user_user_id ON event.event_user(user_id);