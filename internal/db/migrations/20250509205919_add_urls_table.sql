-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS urls (
    user_id UUID, 
    short_url VARCHAR(7), 
    original_url TEXT UNIQUE, 
    is_deleted BOOL DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS urls;
-- +goose StatementEnd
