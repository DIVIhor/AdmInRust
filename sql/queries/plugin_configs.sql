-- name: AddPluginConfig :one
INSERT INTO plugin_configs(plugin_id, config_json, created_at, updated_at)
SELECT p.id, ?, datetime('now'), datetime('now')
FROM plugins AS p
WHERE p.slug = ?
RETURNING *;

-- name: GetPluginConfig :one
SELECT *
FROM plugin_configs
WHERE plugin_id = (
    SELECT id
    FROM plugins
    WHERE slug = ?
);

-- name: UpdatePluginConfig :one
UPDATE plugin_configs
SET config_json = ?, updated_at = datetime('now')
WHERE plugin_id = (
    SELECT id
    FROM plugins
    WHERE slug = ?
)
RETURNING *;

-- name: DeletePluginConfig :one
DELETE
FROM plugin_configs
WHERE plugin_id = (
    SELECT id
    FROM plugins
    WHERE slug = ?
)
RETURNING *;