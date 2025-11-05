CREATE TABLE event_user (
    event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE ,
    user_id UUID not null,

    CREATE INDEX idx_event_user_user_id ON event_user(user_id)
);