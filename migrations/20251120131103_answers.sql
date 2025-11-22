-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE answers(
    id BIGSERIAL PRIMARY KEY NOT NULL,
    question_id BIGINT NOT NULL REFERENCES questions (id) ON DELETE CASCADE,
    user_id UUID NOT NULL DEFAULT gen_random_uuid(),
    answer_text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE answers;
-- +goose StatementEnd
