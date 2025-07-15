-- +goose Up
CREATE TABLE plugin_origins (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    path_to_plugin_list TEXT NOT NULL,
    has_api INTEGER DEFAULT 0 NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

-- +goose Down
DROP TABLE plugin_origins;