-- +goose Up
CREATE TABLE plugins (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    url TEXT NOT NULL,
    origin_id INTEGER NOT NULL,
    is_updated_on_server INTEGER DEFAULT 0 NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,

    FOREIGN KEY (origin_id) REFERENCES plugin_origins(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE plugins;