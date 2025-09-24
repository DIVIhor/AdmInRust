-- name: AddPluginChangelog :one
INSERT INTO plugin_changelogs(plugin_id, version, changelog, update_date, created_at, updated_at)
VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
RETURNING *;

-- name: GetPluginChangelog :many
SELECT *
FROM plugin_changelogs
WHERE plugin_id = (
    SELECT id
    FROM plugins
    WHERE slug = ?
)
ORDER BY update_date DESC, updated_at DESC;