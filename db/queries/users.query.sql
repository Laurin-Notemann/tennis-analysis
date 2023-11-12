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
SELECT *
FROM users
WHERE id = $1 
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1 
LIMIT 1;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1 
LIMIT 1;

-- name: DeleteUserById :one
DELETE FROM users
WHERE id = $1 
RETURNING *;

-- name: UpdateUserById :one
UPDATE users
SET 
  username = $1, 
  email = $2,
  password_hash= $3,
  updated_at = Now()
WHERE id = $4
RETURNING *;

