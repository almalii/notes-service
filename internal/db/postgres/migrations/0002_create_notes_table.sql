-- +goose Up
CREATE TABLE IF NOT EXISTS notes
(
    id         UUID PRIMARY KEY,
    title      TEXT      NOT NULL,
    body       TEXT      NOT NULL,
    tags       TEXT[],
    author     UUID      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);