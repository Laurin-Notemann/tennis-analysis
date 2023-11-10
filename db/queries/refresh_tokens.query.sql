-- name: CreateToken :one
WITH new_token AS (
    INSERT INTO refresh_tokens (
        user_id,
        token,
        expiry_date
    ) VALUES (
        sqlc.arg(user_id),  
        $1,  
        $2
    )
    RETURNING id
)
UPDATE users
SET refresh_token_id = (SELECT id FROM new_token)
WHERE users.id = sqlc.arg(user_id)
RETURNING *;

-- name: GetTokenByUserId :one
SELECT * 
FROM refresh_tokens
WHERE user_id = $1
LIMIT 1;

-- name: UpdateTokenByUserId :one
UPDATE refresh_tokens
SET
  token = $1,
  expiry_date = $2,
  updated_at = Now()
WHERE user_id = $3
RETURNING *;

-- name: DeleteTokenById :exec
DELETE FROM refresh_tokens
WHERE user_id = $1;
