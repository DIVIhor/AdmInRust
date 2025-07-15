-- +goose Up
CREATE TABLE plugin_manual_code_changes (
    id INTEGER PRIMARY KEY,
    plugin_id INTEGER NOT NULL,
    modified_part TEXT NOT NULL,
    row_start INTEGER NOT NULL,
    row_end INTEGER NOT NULL,
    comment TEXT NOT NULL,
    is_relevant INTEGER DEFAULT 0 NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,

    FOREIGN KEY (plugin_id) REFERENCES plugin(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE code_changes;