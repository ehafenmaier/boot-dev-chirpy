-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (gen_random_uuid(), now(), now(), $1, $2)
RETURNING *;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password
FROM users
WHERE email = $1
LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET updated_at = now(),
    email = $1,
    hashed_password = $2
WHERE id = $3
RETURNING *;