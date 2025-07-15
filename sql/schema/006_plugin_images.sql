-- +goose Up
CREATE TABLE plugin_images (
    id INTEGER PRIMARY KEY,
    plugin_id INTEGER NOT NULL,
    image BLOB NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,

    FOREIGN KEY (plugin_id) REFERENCES plugins(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE plugin_images;