CREATE TABLE course (
    id UUID PRIMARY KEY,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL
);

-- name: CreateCourse :one
INSERT INTO course (id, code, name, description) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetCourse :one
SELECT * FROM course
WHERE id = $1;

-- name: UpdateCourse :one
UPDATE course
SET code = $2,
    name = $3,
    description = $4
WHERE id = $1
RETURNING *;

-- name: DeleteCourse :exec
DELETE FROM course
WHERE id = $1;