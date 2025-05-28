CREATE TABLE user_courses (
    user_id INT NOT NULL,              
    course_id INT NOT NULL,             
    enrollment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  
    PRIMARY KEY (user_id, course_id),   
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,  
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE 
);