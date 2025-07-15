-- +goose Up
CREATE TABLE plugin_versions (
    id INTEGER PRIMARY KEY,
    plugin_id INTEGER NOT NULL,
    version TEXT NOT NULL,
    change_log TEXT NOT NULL,
    update_date TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,

    FOREIGN KEY (plugin_id) REFERENCES plugins(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE plugin_versions;