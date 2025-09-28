-- name: CreateNode :one
INSERT INTO nodes (
    namespace, name, image, cmd, container_id
) VALUES (
    ?, ?, ?, ?, ?
)
ON CONFLICT DO UPDATE SET
    image = excluded.image,
    cmd = excluded.cmd,
    container_id = excluded.container_id,
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

-- name: UpdateNode :exec
UPDATE nodes
SET image = ?, cmd = ?, container_id = ?, updated_at = CURRENT_TIMESTAMP
WHERE namespace = ? AND name = ? AND deleted_at IS NULL;

-- name: SoftDeleteNode :exec
UPDATE nodes
SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE namespace = ? AND name = ? AND deleted_at IS NULL;