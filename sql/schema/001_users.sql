-- +goose up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    name TEXT NOT NULL,
    favorite_color TEXT,
    UNIQUE (name)
);

-- +goose down
DROP TABLE users;