-- +goose Up
CREATE TABLE plugin_locales (
    plugin_id INTEGER NOT NULL,
    lang_code TEXT NOT NULL,
    lang_name TEXT NOT NULL,
    content_json TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,

    PRIMARY KEY (plugin_id, lang_code),
    FOREIGN KEY (plugin_id) REFERENCES plugins(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE plugin_locales;