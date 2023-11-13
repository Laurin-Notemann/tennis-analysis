-- name: CreateNewTeamWithOnePlayer :one
WITH new_player AS (
  INSERT INTO players (
    first_name,
    last_name
  ) VALUES (
    $1,
    $2
  )
  RETURNING id
)
INSERT INTO teams (
  name,
  user_id,
  player_one
) VALUES (
  $3,
  $4,
  (SELECT id FROM new_player)
)
RETURNING *;

-- name: CreateTeamWithTwoPlayers :one
INSERT INTO teams (
  name,
  user_id,
  player_one,
  player_two
) VALUES (
  $1,
  $2,
  $3,
  $4
)
RETURNING *;

-- name: GetTeamById :one
SELECT * 
FROM teams
WHERE id = $1
LIMIT 1;

-- name: DeleteTeamById :one
DELETE FROM teams
WHERE id = $1
RETURNING *;

