-- +goose Up
-- +goose StatementBegin
CREATE TABLE questions (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    question_text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE questions;
-- +goose StatementEnd
