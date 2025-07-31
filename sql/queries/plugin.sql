-- name: AddPlugin :one
INSERT INTO plugins(name, description, url, origin_id, is_updated_on_server, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, datetime('now'), datetime('now'))
RETURNING *;

-- name: GetPlugins :many
SELECT *
FROM plugins
ORDER BY is_updated_on_server, updated_at, origin_id, name;

-- name: GetPlugin :one
SELECT *
FROM plugins
WHERE id = ?;

-- name: DeletePlugin :one
DELETE
FROM plugins
WHERE id = ?
RETURNING *;