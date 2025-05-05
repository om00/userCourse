package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/om00/userCourse/models"
	"github.com/om00/userCourse/mysqldb"
)

type App struct {
	Db        *mysqldb.DbIns
	Validator *validator.Validate
}

func (app *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := app.Validator.Struct(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = app.Db.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
	})
}

func (app *App) GetAllCourse(w http.ResponseWriter, r *http.Request) {
	var id int64
	var err error
	idStr := r.URL.Query().Get("id")
	name := r.URL.Query().Get("name")

	if idStr != "" {
		id, err = strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user_id", http.StatusBadRequest)
			return
		}
	}

	courses, err := app.Db.GetAllCourse(id, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(courses); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (app *App) UserCourse(w http.ResponseWriter, r *http.Request) {
	var req models.UserCourseReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := app.Validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validActions := map[string]bool{
		"assign":   true,
		"unassign": true,
	}

	if !validActions[req.Action] {
		http.Error(w, "Invalid action. Must be 'assign' or 'unassign'", http.StatusBadRequest)
		return
	}

	user, err := app.Db.GetUser(int64(req.UserID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	course, err := app.Db.GetCourse(req.CourseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch req.Action {
	case "assign":
		_, err = app.Db.CreateUserCourse(user.ID, course.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Course is assign successFully",
		})

	case "unassign":
		_, err = app.Db.DeleteUserCourse(user.ID, course.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Course is Un-assign successFully",
		})

	}

}

func (app *App) UserDetials(w http.ResponseWriter, r *http.Request) {
	var id int64
	var err error
	idStr := r.URL.Query().Get("id")
	name := r.URL.Query().Get("name")

	if idStr != "" {
		id, err = strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user_id", http.StatusBadRequest)
			return
		}
	}

	userMap, err := app.Db.GetUserDetialsMap(id, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var userIds []int64
	for k := range userMap {
		userIds = append(userIds, k)
	}

	userCourses, err := app.Db.GetUserCourse(userIds)
	if err != nil {

	}

	var userInfo []models.User
	for _, v := range userCourses {
		if user, ok := userMap[v.UserID]; ok {
			courseInfo := models.UserCourse{UserID: v.UserID, CourseID: v.CourseID, EnrolledAt: v.EnrolledAt}
			user.UserCourses = append(user.UserCourses, courseInfo)
			userMap[v.UserID] = user
		}
	}

	for _, v := range userMap {
		userInfo = append(userInfo, v)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(userInfo); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
