CREATE TABLE enrollment (
    id UUID PRIMARY KEY,
    course_offering_id UUID REFERENCES course_offering(id),
    student_id UUID REFERENCES student(id),
    enrollment_date DATE NOT NULL
);