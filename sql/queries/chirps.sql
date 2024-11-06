-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (gen_random_uuid(), now(), now(), $1, $2)
RETURNING *;

-- name: GetChirps :many
SELECT id, created_at, updated_at, body, user_id
FROM chirps
ORDER BY CASE WHEN $1 = 'desc' THEN created_at END DESC,
         CASE WHEN $1 = 'asc' THEN created_at END ASC;

-- name: GetChirpsByUserId :many
SELECT id, created_at, updated_at, body, user_id
FROM chirps
WHERE user_id = $1
ORDER BY CASE WHEN $2 = 'desc' THEN created_at END DESC,
         CASE WHEN $2 = 'asc' THEN created_at END ASC;

-- name: GetChirpById :one
SELECT id, created_at, updated_at, body, user_id
FROM chirps
WHERE id = $1
LIMIT 1;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;