-- name: CreateNode :one
INSERT INTO nodes (
    namespace, name, image, cmd, container_id, resource_version
) VALUES (
    ?, ?, ?, ?, ?, 1
)
ON CONFLICT DO UPDATE SET
    image = excluded.image,
    cmd = excluded.cmd,
    container_id = excluded.container_id,
    resource_version = excluded.resource_version + 1,
    deleted_at = NULL,
    updated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: GetNode :one
SELECT * FROM nodes
WHERE namespace = ? AND name = ? AND deleted_at IS NULL;

-- name: GetNodeWithDeleted :one
SELECT * FROM nodes
WHERE namespace = ? AND name = ?;

-- name: ListNodes :many
SELECT * FROM nodes
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListNodesWithDeleted :many
SELECT * FROM nodes
ORDER BY created_at DESC;

-- name: UpdateNode :one
UPDATE nodes
SET image = ?, cmd = ?, container_id = ?, updated_at = CURRENT_TIMESTAMP,  resource_version = resource_version + 1
WHERE namespace = ? AND name = ? AND deleted_at IS NULL AND resource_version = ?
RETURNING *;

-- name: SoftDeleteNode :one
UPDATE nodes
SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP, resource_version = resource_version + 1
WHERE namespace = ? AND name = ? AND deleted_at IS NULL AND resource_version = ?
RETURNING *;