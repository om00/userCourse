package seeders

import (
	"database/sql"
	"fmt"
	"log"
)

func SeedCourses(db *sql.DB) {
	_, err := db.Exec(`
        INSERT INTO courses (name, description, created_at)
VALUES
    ('Introduction to Go Programming', 'A beginner-friendly course to learn Go programming language.', now()),
    ('Advanced JavaScript', 'An in-depth course focusing on advanced JavaScript concepts and techniques.', now()),
    ('Database Management Systems', 'A comprehensive course covering relational and non-relational databases.', now()),
    ('Web Development with React', 'Learn modern web development using React and JavaScript.', now()),
    ('Data Science with Python', 'An introductory course on data science using Python and its libraries.', now()),
    ('Machine Learning Basics', 'A fundamental course to understand the basics of machine learning algorithms.', now());
    `)
	if err != nil {
		log.Fatalf("Error seeding course table: %v", err)
	}
	fmt.Println("course  data seeded successfully")
}
