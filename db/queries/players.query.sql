-- name: GetPlayerById :one
SELECT * 
FROM players
WHERE id = $1
LIMIT 1;

-- name: DeletePlayerById :one
DELETE FROM players
WHERE id = $1
RETURNING *;

-- name: UpdatePlayerById :one
UPDATE players
SET 
  first_name = $1, 
  last_name = $2,
  updated_at = Now()
WHERE id = $3
RETURNING *;
