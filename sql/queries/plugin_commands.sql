-- name: AddPluginCommand :one
INSERT INTO plugin_commands(plugin_id, command, description, created_at, updated_at)
VALUES (?, ?, ?, datetime('now'), datetime('now'))
RETURNING *;

-- name: GetPluginCommands :many
SELECT *
FROM plugin_commands
WHERE plugin_id = (
    SELECT id
    FROM plugins
    WHERE slug = ?
);

-- name: DeletePluginCommand :one
DELETE
FROM plugin_commands
WHERE plugin_id = ?
RETURNING *;