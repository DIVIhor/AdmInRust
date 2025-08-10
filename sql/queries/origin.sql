-- name: AddOrigin :one
INSERT INTO plugin_origins(name, slug, url, path_to_plugin_list, has_api, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, datetime('now'), datetime('now'))
RETURNING *;

-- name: GetOrigins :many
SELECT *
FROM plugin_origins
ORDER BY name;

-- name: GetOrigin :one
SELECT *
FROM plugin_origins
WHERE slug = ?;

-- name: DeleteOrigin :one
DELETE
FROM plugin_origins
WHERE slug = ?
RETURNING *;