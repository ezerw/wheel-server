-- name: ListPeople :many
SELECT id, first_name, last_name, email, team_id
FROM people
WHERE team_id = ?
ORDER BY id;

-- name: GetPerson :one
SELECT id, first_name, last_name, email, team_id
FROM people
WHERE id = ?
  AND team_id = ?
LIMIT 1;

-- name: CreatePerson :execresult
INSERT INTO people (
    first_name,
    last_name,
    email,
    team_id
) VALUES (
    ?, ?, ?, ?
);

-- name: UpdatePerson :execresult
UPDATE people
SET first_name = ?, last_name = ?, email = ?, team_id = ?
WHERE id = ?;

-- name: DeletePerson :exec
DELETE FROM people
WHERE id = ?
AND team_id = ?;