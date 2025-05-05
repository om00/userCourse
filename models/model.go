package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Age       int       `json:"age" validate:"required,min=1"`
	CreatedAt time.Time `json:"-"`

	UserCourses []UserCourse `json:"user_courses"`
}

type UserCourse struct {
	UserID     int64     `json:"user_id"`
	CourseID   int64     `json:"course_id"`
	EnrolledAt time.Time `json:"enrolled_at"`
}

type Course struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserCourseReq struct {
	UserID   int64  `json:"user_id"`
	CourseID int64  `json:"course_id"`
	Action   string `json:"action"`
}
