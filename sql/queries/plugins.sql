-- name: AddPlugin :one
INSERT INTO plugins(name, slug, description, url, origin_id, is_updated_on_server, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))
RETURNING *;

-- name: GetPlugins :many
SELECT *
FROM plugins
ORDER BY is_updated_on_server, updated_at, origin_id, name;

-- name: GetPlugin :one
SELECT *
FROM plugins
WHERE slug = ?;

-- name: GetPluginWithOriginsJson :one
SELECT plugins.*, (
    SELECT json_group_array(json_object('id', origins.id, 'name', origins.name))
    FROM plugin_origins origins
) AS origin_options
FROM plugins
WHERE plugins.slug = ?;

-- name: UpdatePlugin :one
UPDATE plugins
SET description = ?,
    url = ?,
    origin_id = ?,
    is_updated_on_server = ?,
    updated_at = datetime('now')
WHERE slug = ?
RETURNING *;

-- name: DeletePlugin :one
DELETE
FROM plugins
WHERE slug = ?
RETURNING *;