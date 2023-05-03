CREATE TABLE course_offering (
    id UUID PRIMARY KEY,
    course_id UUID REFERENCES course(id),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL
);