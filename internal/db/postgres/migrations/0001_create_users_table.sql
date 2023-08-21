-- +goose Up
CREATE TABLE IF NOT EXISTS users
(
    id         UUID PRIMARY KEY,
    username   TEXT      NOT NULL,
    email      TEXT      NOT NULL UNIQUE,
    password   TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);


