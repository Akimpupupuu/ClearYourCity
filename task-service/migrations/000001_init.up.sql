CREATE SCHEMA IF NOT EXISTS task_service;

CREATE TABLE IF NOT EXISTS task_service.task (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    user_id BIGINT NOT NULL,
    title varchar(100) NOT NULL CHECK(char_length(title) BETWEEN 3 AND 100),
    description varchar(1000) NOT NULL CHECK(char_length(description) BETWEEN 10 AND 1000),
    status varchar(100) NOT NULL DEFAULT 'created' CHECK(status IN ('created', 'in_progress', 'done', 'rejected')),
    created_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ,

    CHECK(
        (status IN ('created', 'in_progress') AND completed_at IS NULL)
        OR
        (status IN ('done', 'rejected') AND completed_at IS NOT NULL AND completed_at >= created_at)
    )
);