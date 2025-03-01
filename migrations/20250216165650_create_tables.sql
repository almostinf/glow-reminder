-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS reminders (
    id UUID NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    msg TEXT NOT NULL,
    colour SMALLINT NOT NULL,
    mode SMALLINT NOT NULL,
    scheduled_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reminders;
-- +goose StatementEnd