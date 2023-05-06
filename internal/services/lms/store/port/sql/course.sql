CREATE TABLE course (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL
);

-- name: CreateCourse :exec
INSERT INTO course (id, name, description) VALUES ($1, $2, $3);

-- name: GetCourseById :one
SELECT * FROM course WHERE id = $1;

-- name: UpdateCourseById :exec
UPDATE course SET name = $2, description = $3 WHERE id = $1;

-- name: UpdateCourseNameById :exec
UPDATE course SET name = $2 WHERE id = $1;

-- name: UpdateCourseDescriptionById :exec
UPDATE course SET description = $2 WHERE id = $1;

-- name: DeleteCourse :exec
DELETE FROM course WHERE id = $1;