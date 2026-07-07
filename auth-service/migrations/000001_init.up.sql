CREATE SCHEMA IF NOT EXISTS auth_service;

CREATE TABLE IF NOT EXISTS auth_service.users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    full_name VARCHAR(100) NOT NULL CHECK(char_length(full_name) BETWEEN 3 AND 100),
    email VARCHAR(100) NOT NULL UNIQUE CHECK(
        email ~ '^[^@\s]+@[^@\s]+\.[^@\s]+$'
        AND
        char_length(email) BETWEEN 5 AND 100
    ),
    hashed_password TEXT NOT NULL,
    created_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS auth_service.sessions (
    id UUID PRIMARY KEY NOT NULL,
    user_id BIGINT NOT NULL REFERENCES auth_service.users(id) ON DELETE CASCADE,
    refresh_token TEXT NOT NULL,
    is_revoked BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL
);