CREATE TYPE event_status AS ENUM (
    'draft',
    'published',
    'ongoing',
    'completed',
    'cancelled',
    'postponed'
);

CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    about TEXT,
    start_date TIMESTAMP NOT NULL,
    location VARCHAR(255),
    status event_status NOT NULL,
    max_attendees SMALLINT,
    current_attendance SMALLINT DEFAULT 0,
    creator UUID NOT NULL
);