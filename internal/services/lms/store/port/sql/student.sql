CREATE TABLE student (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

-- name: CreateStudent :exec
INSERT INTO student (id, name, email, password) VALUES ($1, $2, $3, $4);

-- name: GetStudentById :one
SELECT * FROM student WHERE id = $1;

-- name: UpdateStudentById :exec
UPDATE student SET name = $2 WHERE id = $1;

-- name: DeleteStudent :execrows
DELETE FROM student WHERE id = $1;