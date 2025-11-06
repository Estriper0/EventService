CREATE TABLE event_user (
    event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    UNIQUE(event_id, user_id)
);

CREATE INDEX idx_event_user_user_id ON event_user(user_id);