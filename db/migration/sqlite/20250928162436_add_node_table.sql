-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS nodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    -- metadata
    namespace TEXT NOT NULL,
    name TEXT NOT NULL,
    -- spec
    image TEXT NOT NULL,
    cmd TEXT NOT NULL,
    -- status
    container_id TEXT,
    -- misc
    deleted_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(namespace, name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS nodes;
-- +goose StatementEnd