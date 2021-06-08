-- name: GetTeam :one
SELECT id, name
FROM teams
WHERE id = ?
LIMIT 1;

-- name: ListTeams :many
SELECT id, name
FROM teams
ORDER BY id;

-- name: CreateTeam :execresult
INSERT INTO teams (name)
VALUES ( ? );

-- name: UpdateTeam :execresult
UPDATE teams
SET name = ?
WHERE id = ?;

-- name: DeleteTeam :exec
DELETE FROM teams
WHERE id = ?;