-- name: ListTurns :many
SELECT t.id, t.person_id, t.date, t.created_at
FROM turns t
         LEFT JOIN people p ON t.person_id = p.id
WHERE p.team_id = ?
ORDER BY t.date DESC
LIMIT ? OFFSET ?;

-- name: ListTurnsWithDateFrom :many
SELECT t.id, t.person_id, t.date, t.created_at
FROM turns t
         LEFT JOIN people p ON t.person_id = p.id
WHERE p.team_id = ?
  AND t.date >= ?
ORDER BY t.date DESC
LIMIT ? OFFSET ?;

-- name: ListTurnsWithDateTo :many
SELECT t.id, t.person_id, t.date, t.created_at
FROM turns t
         LEFT JOIN people p ON t.person_id = p.id
WHERE p.team_id = ?
  AND t.date <= ?
ORDER BY t.date DESC
LIMIT ? OFFSET ?;

-- name: ListTurnsWithBothDates :many
SELECT t.id, t.person_id, t.date, t.created_at
FROM turns t
         LEFT JOIN people p ON t.person_id = p.id
WHERE p.team_id = ?
  AND t.date >= ?
  AND t.date <= ?
ORDER BY t.date DESC
LIMIT ? OFFSET ?;

-- name: GetTurn :one
SELECT t.id, t.person_id, t.date, t.created_at
FROM turns t
         LEFT JOIN people p ON t.person_id = p.id
WHERE t.id = ?
  AND p.team_id = ?
LIMIT 1;

-- name: GetTurnByDate :one
SELECT t.id, t.person_id, t.date, t.created_at
FROM turns t
         LEFT JOIN people p ON t.person_id = p.id
WHERE t.date = ?
  AND p.team_id = ?
LIMIT 1;

-- name: GetTurnByDateAndTeam :one
SELECT t.id, t.person_id, t.date, t.created_at
FROM turns t
         LEFT JOIN people p ON t.person_id = p.id
WHERE t.date = ?
  AND p.team_id = ?
LIMIT 1;

-- name: CreateTurn :execresult
INSERT INTO turns (person_id, date)
VALUES (?, ?);

-- name: UpdateTurn :execresult
UPDATE turns
SET person_id = ?,
    date      = ?
WHERE id = ?;

-- name: DeleteTurn :exec
DELETE
FROM turns
WHERE id = ?
  AND person_id = ?;