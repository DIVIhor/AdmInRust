-- name: AddPluginDoc :one
INSERT INTO plugin_docs(plugin_id, doc, created_at, updated_at)
VALUES (?, ?, datetime('now'), datetime('now'))
RETURNING *;

-- name: GetPluginDoc :one
SELECT *
FROM plugin_docs
WHERE plugin_id = (
    SELECT id
    FROM plugins
    WHERE slug = ?
);

-- name: UpdatePluginDoc :one
UPDATE plugin_docs
SET doc = ?,
    updated_at = datetime('now')
WHERE plugin_id = (
    SELECT id
    FROM plugins
    WHERE slug = ?
)
RETURNING *;

-- name: DeletePluginDoc :one
DELETE
FROM plugin_docs
WHERE plugin_id = (
    SELECT id
    FROM plugins
    WHERE slug = ?
)
RETURNING *;