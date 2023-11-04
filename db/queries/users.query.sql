-- name: CreateUser :one
INSERT INTO users (
  username,
  email,
  password_hash
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetAllUsers :many
SELECT * 
FROM users;

-- name: GetUserById :one
SELECT id, username, email 
FROM users
WHERE id = $1 
LIMIT 1;
