CREATE TABLE IF NOT EXISTS student (
  id SERIAL PRIMARY KEY,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  phone VARCHAR(20),
  address VARCHAR(255),
  date_of_birth DATE NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- name: GetStudent :one
SELECT * FROM student WHERE id = $1;

-- name: ListStudents :many
SELECT * FROM student ORDER BY last_name, first_name;

-- name: CreateStudent :exec
INSERT INTO student (first_name, last_name, email, phone, address, date_of_birth)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: UpdateStudent :exec
UPDATE student
SET first_name = $2, last_name = $3, email = $4, phone = $5, address = $6, date_of_birth = $7, updated_at = NOW()
WHERE id = $1;

-- name: DeleteStudent :exec
DELETE FROM student WHERE id = $1;

-- name: EnrollStudent :exec
INSERT INTO enrollments (student_id, course_code, semester, grade)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: GetStudentEnrollments :many
SELECT student.*, enrollments.course_code, enrollments.semester, enrollments.grade
FROM student
JOIN enrollments ON student.id = enrollments.student_id
WHERE student.id = $1;