CREATE TABLE IF NOT EXISTS enrollments (
  id SERIAL PRIMARY KEY,
  student_id INT NOT NULL REFERENCES students(id),
  course_code VARCHAR(20) NOT NULL,
  semester VARCHAR(10) NOT NULL,
  grade VARCHAR(2),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- name: EnrollStudent :exec
INSERT INTO enrollments (student_id, course_code, semester, grade)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: GetStudentEnrollments :many
SELECT student.*, enrollments.course_code, enrollments.semester, enrollments.grade
FROM student
JOIN enrollments ON student.id = enrollments.student_id
WHERE student.id = $1;