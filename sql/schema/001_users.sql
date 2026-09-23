-- +goose up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    
    name TEXT NOT NULL,
    favorite_color TEXT,
    UNIQUE (name)
);

-- +goose down
DROP TABLE users;