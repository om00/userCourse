CREATE TABLE user_courses (
    user_id INT NOT NULL,               -- Foreign key referencing the users table
    course_id INT NOT NULL,             -- Foreign key referencing the courses table
    enrollment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Date and time when the user was enrolled
    PRIMARY KEY (user_id, course_id),   -- Composite primary key to ensure a user can't enroll in the same course multiple times
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,  -- Ensures referential integrity
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE  -- Ensures referential integrity
);