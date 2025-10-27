-- name: AddPluginLocale :one
INSERT INTO plugin_locales(
    plugin_id, lang_code,
    lang_name, content_json,
    created_at, updated_at
)
SELECT p.id, ?, ?, ?, datetime('now'), datetime('now')
FROM plugins AS p
WHERE p.slug = ?
RETURNING *;

-- name: GetPluginLocales :many
SELECT *
FROM plugin_locales
WHERE plugin_id = (
    SELECT id
    FROM plugins
    WHERE slug = ?
);

-- name: GetPluginLocale :one
SELECT *
FROM plugin_locales
WHERE plugin_id = (
    SELECT id
    FROM plugins
    WHERE slug = ?
) AND lang_code = ?;

-- name: UpdatePluginLocale :one
UPDATE plugin_locales
SET content_json = ?,
    updated_at = datetime('now')
WHERE plugin_id = (
    SELECT id
    FROM plugins
    WHERE slug = ?
) AND lang_code = ?
RETURNING *;

-- name: DeletePluginLocale :one
DELETE
FROM plugin_locales
WHERE lang_code = ? AND plugin_id = (
    SELECT id
    FROM plugins
    WHERE slug = ?
)
RETURNING *;