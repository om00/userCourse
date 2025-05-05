package mysqldb

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/om00/userCourse/models"
	seeders "github.com/om00/userCourse/mysqldb/seeders"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var Dbpath string

type DbIns struct {
	mainDB *sql.DB
}

func RunMigrations(cmd string) {
	migrationsDir := "file://mysqldb/migrations"

	m, err := migrate.New(
		migrationsDir,
		Dbpath,
	)
	if err != nil {
		log.Fatalf("Could not initialize migration: %v", err)
	}

	defer m.Close()

	switch cmd {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration up failed: %v", err)
		}
		fmt.Println("Migration UP done successfully.")

	case "down":
		if err := m.Down(); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
		fmt.Println("Migration DOWN done successfully.")

	case "drop":
		if err := m.Drop(); err != nil {
			log.Fatalf("Drop failed: %v", err)
		}
		fmt.Println("Database dropped successfully.")

	default:
		fmt.Println("Unknown migration command:", cmd)
	}
}

func CallSeederFunction(db *sql.DB, funcName string) {
	seedersMap := getSeederFunctions()

	fmt.Println(db)
	if seedFunc, exists := seedersMap[funcName]; exists {

		reflect.ValueOf(seedFunc).Call([]reflect.Value{reflect.ValueOf(db)})

	} else {
		log.Printf("Seeder function %s not found.", funcName)
	}
}

func getSeederFunctions() map[string]interface{} {
	seedersMap := make(map[string]interface{})

	// Add the seeder functions from the seeders package to the map
	seedersMap["seedUsers"] = seeders.SeedUsers
	seedersMap["seedCoures"] = seeders.SeedCourses

	return seedersMap
}

func NewDB(db *sql.DB) *DbIns {
	return &DbIns{mainDB: db}
}

func (db *DbIns) CreateUser(user models.User) (int64, error) {
	query := `INSERT INTO users (name, email, age, created_at) VALUES (?, ?, ?, ?)`

	result, err := db.mainDB.Exec(query, user.Name, user.Email, user.Age, time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return lastInsertID, nil
}

func (db *DbIns) GetAllCourse(id int64, name string) ([]models.Course, error) {
	query := `
	SELECT c.id, c.name, c.description
	FROM courses c
`
	args := []interface{}{}
	conditions := []string{}

	if id > 0 {
		conditions = append(conditions, "c.id = ?")
		args = append(args, id)
	}

	if name != "" {
		conditions = append(conditions, "c.name LIKE ?")
		args = append(args, "%"+name+"%")
	}

	if len(conditions) > 0 {
		query += "WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := db.mainDB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch courses: %w", err)
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		if err := rows.Scan(&course.ID, &course.Name, &course.Description); err != nil {
			return nil, fmt.Errorf("failed to scan course: %w", err)
		}
		courses = append(courses, course)
	}

	return courses, nil
}

func (db *DbIns) GetUser(id int64) (models.User, error) {

	query := "SELECT id, name, email, age, created_at FROM users WHERE id = ?"
	var user models.User

	err := db.mainDB.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.CreatedAt)
	if err != nil {

		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("user with id %d not found", id)
		}
		return models.User{}, fmt.Errorf("failed to fetch user: %w", err)
	}

	return user, nil
}

func (db *DbIns) GetCourse(id int64) (models.Course, error) {

	query := "SELECT id, name, description, created_at FROM courses WHERE id = ?"
	var course models.Course

	err := db.mainDB.QueryRow(query, id).Scan(&course.ID, &course.Name, &course.Description, &course.CreatedAt)
	if err != nil {

		if err == sql.ErrNoRows {
			return models.Course{}, fmt.Errorf("course with id %d not found", id)
		}
		return models.Course{}, fmt.Errorf("failed to fetch course: %w", err)
	}

	return course, nil
}

func (db *DbIns) CreateUserCourse(userId, courseId int64) (int64, error) {

	query := `INSERT INTO user_courses (user_id, course_id, enrollment_date) 
			  VALUES (?, ?, NOW())`

	result, err := db.mainDB.Exec(query, userId, courseId)
	if err != nil {
		return 0, fmt.Errorf("failed to assign course to user: %w", err)
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	return insertedID, nil
}

func (db *DbIns) DeleteUserCourse(userId, courseId int64) (int64, error) {

	query := `DELETE FROM user_courses WHERE user_id = ? AND course_id = ?`

	result, err := db.mainDB.Exec(query, userId, courseId)
	if err != nil {
		return 0, fmt.Errorf("failed to delete user-course association: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to check affected rows: %w", err)
	}

	return rowsAffected, nil
}

func (db *DbIns) GetUserDetialsMap(id int64, name string) (map[int64]models.User, error) {

	usersMap := make(map[int64]models.User)

	query := "SELECT id, name, email, age, created_at FROM users WHERE 1=1"
	args := []interface{}{}

	if id > 0 {
		query += " AND id = ?"
		args = append(args, id)
	}
	if name != "" {
		query += " AND name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	rows, err := db.mainDB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		usersMap[user.ID] = user
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return usersMap, nil

}

func (db *DbIns) GetUserCourse(ids []int64) ([]models.UserCourse, error) {
	var courses []models.UserCourse

	if len(ids) == 0 {
		return courses, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(
		`SELECT user_id, course_id, enrollment_date FROM user_courses WHERE user_id IN (%s)`,
		strings.Join(placeholders, ","),
	)

	rows, err := db.mainDB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var uc models.UserCourse
		if err := rows.Scan(&uc.UserID, &uc.CourseID, &uc.EnrolledAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		courses = append(courses, uc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return courses, nil
}
