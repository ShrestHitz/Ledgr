-- name: CreateCategory :one
INSERT INTO categories (
  owner,
  name,
  type
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetCategory :one
SELECT * FROM categories
WHERE id = $1 LIMIT 1;

-- name: ListCategories :many
SELECT * FROM categories
WHERE owner IS NULL OR owner = $1
ORDER BY name;

-- name: ListCategoriesByType :many
SELECT * FROM categories
WHERE (owner IS NULL OR owner = $1) AND type = $2
ORDER BY name;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;
