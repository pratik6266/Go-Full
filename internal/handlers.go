package internal

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewHandler(d *Db, logger *logrus.Logger) *Handler {
	return &Handler{
		db:     d.db,
		logger: logger,
	}
}

// Healthcheck godoc
// @Summary      Health check
// @Description  Returns the API health status.
// @Tags         Health
// @ID           healthcheck
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /health [get]
func (h *Handler) Healthcheck(c *gin.Context) {
	h.logger.Info("Healthcheck endpoint called")
	c.JSON(http.StatusOK, HealthResponse{Message: "API is healthy"})
}

// GetStudents godoc
// @Summary      List students
// @Description  Returns all students
// @Tags         Students
// @Produce      json
// @Success      200  {array}   Student
// @Failure      500  {object}  ErrorResponse
// @Router       /students [get]
func (h *Handler) GetStudents(c *gin.Context) {
	rows, err := h.db.Query("SELECT id, name, age, email FROM students")
	if err != nil {
		h.logger.Error("Error querying students:", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch students"})
		return
	}
	defer rows.Close()

	students := make([]Student, 0)
	for rows.Next() {
		var s Student
		if err := rows.Scan(&s.ID, &s.Name, &s.Age, &s.Email); err != nil {
			h.logger.Error("Error scanning student row:", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to parse student record"})
			return
		}
		students = append(students, s)
	}
	if err := rows.Err(); err != nil {
		h.logger.Error("Row iteration error:", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database error while fetching students"})
		return
	}

	h.logger.Info("Fetched students successfully")
	c.JSON(http.StatusOK, students)
}

// CreateStudent godoc
// @Summary      Create a student
// @Description  Creates a new student
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        student  body      StudentCreateRequest  true  "Student payload"
// @Success      201      {object}  Student
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /students [post]
func (h *Handler) CreateStudent(c *gin.Context) {
	var req StudentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request payload:", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	var student Student
	err := h.db.QueryRow(
		"INSERT INTO students (name, age, email) VALUES ($1, $2, $3) RETURNING id",
		req.Name, req.Age, req.Email,
	).Scan(&student.ID)
	if err != nil {
		h.logger.Error("Failed to create student:", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create student"})
		return
	}
	student.Name = req.Name
	student.Age = req.Age
	student.Email = req.Email

	h.logger.Info("Created student successfully with ID:", student.ID)
	c.JSON(http.StatusCreated, student)
}

// GetStudentByID godoc
// @Summary      Get a student by ID
// @Description  Retrieves a student by its ID
// @Tags         Students
// @Produce      json
// @Param        id   path      int  true  "Student ID"
// @Success      200  {object}  Student
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /students/{id} [get]
func (h *Handler) GetStudentByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.logger.Error("Invalid student ID:", idStr)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	var s Student
	err = h.db.QueryRow("SELECT id, name, age, email FROM students WHERE id = $1", id).
		Scan(&s.ID, &s.Name, &s.Age, &s.Email)
	if err == sql.ErrNoRows {
		h.logger.Warn("Student not found with ID:", id)
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Student not found"})
		return
	} else if err != nil {
		h.logger.Error("Failed to fetch student:", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch student"})
		return
	}

	h.logger.Info("Fetched student successfully with ID:", id)
	c.JSON(http.StatusOK, s)
}

// UpdateStudent godoc
// @Summary      Update a student
// @Description  Updates an existing student by ID
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        id       path      int                      true  "Student ID"
// @Param        student  body      StudentUpdateRequest     true  "Student payload"
// @Success      200      {object}  Student
// @Failure      400      {object}  ErrorResponse
// @Failure      404      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /students/{id} [put]
func (h *Handler) UpdateStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.logger.Error("Invalid student ID for update:", idStr)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	var payload StudentUpdateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.logger.Error("Invalid update payload:", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	// Update and return the updated row
	var updated Student
	err = h.db.QueryRow(
		"UPDATE students SET name = $1, age = $2, email = $3 WHERE id = $4 RETURNING id, name, age, email",
		payload.Name, payload.Age, payload.Email, id,
	).Scan(&updated.ID, &updated.Name, &updated.Age, &updated.Email)
	if err == sql.ErrNoRows {
		h.logger.Warn("Student not found for update with ID:", id)
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Student not found"})
		return
	} else if err != nil {
		h.logger.Error("Failed to update student:", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update student"})
		return
	}

	h.logger.Info("Updated student successfully with ID:", id)
	c.JSON(http.StatusOK, updated)
}

// DeleteStudent godoc
// @Summary      Delete a student
// @Description  Deletes a student by ID
// @Tags         Students
// @Param        id  path  int  true  "Student ID"
// @Success      204
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /students/{id} [delete]
func (h *Handler) DeleteStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.logger.Error("Invalid student ID for deletion:", idStr)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	res, err := h.db.Exec("DELETE FROM students WHERE id = $1", id)
	if err != nil {
		h.logger.Error("Failed to delete student:", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete student"})
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		h.logger.Warn("Student not found for deletion with ID:", id)
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Student not found"})
		return
	}

	h.logger.Info("Deleted student successfully with ID:", id)
	c.Status(http.StatusNoContent)
}

// GetUsers godoc
// @Summary      List users
// @Description  Returns all users
// @Tags         Users
// @Produce      json
// @Success      200  {array}   User
// @Failure      500  {object}  ErrorResponse
// @Router       /users [get]
// GetUsers returns the list of users from the database.
func (h *Handler) GetUsers(c *gin.Context) {
	query := "SELECT id, name, email FROM users"
	rows, err := h.db.Query(query)
	if err != nil {
		h.logger.Error("Error querying users:", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch users"})
		return
	}

	defer rows.Close()
	users := make([]User, 0)
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			h.logger.Error("Error scanning user row:", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to parse user record"})
			return
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		h.logger.Error("Row iteration error:", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database error while fetching users"})
		return
	}

	h.logger.Info("Fetched users successfully")
	c.JSON(http.StatusOK, users)
}

// CreateUser godoc
// @Summary      Create a user
// @Description  Creates a new user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user  body      UserCreateRequest  true  "User payload"
// @Success      201   {object}  User
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var req UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid user creation payload:", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	var user User
	err := h.db.QueryRow(
		"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
		req.Name, req.Email,
	).Scan(&user.ID)
	if err != nil {
		h.logger.Error("Failed to create user:", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create user"})
		return
	}
	user.Name = req.Name
	user.Email = req.Email

	h.logger.Info("Created user successfully with ID:", user.ID)
	c.JSON(http.StatusCreated, user)
}

// GetUserById godoc
// @Summary      Get user by ID
// @Description  Retrieve a single user by ID using query parameter ?id=
// @Tags         Users
// @Produce      json
// @Param        id   query     int  true  "User ID"
// @Success      200  {object}  User
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /users/by-id [get]
func (h *Handler) GetUserById(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.logger.Error("Invalid user ID:", idStr)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	var u User
	err = h.db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id).
		Scan(&u.ID, &u.Name, &u.Email)
	if err == sql.ErrNoRows {
		h.logger.Warn("User not found with ID:", id)
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		return
	} else if err != nil {
		h.logger.Error("Failed to fetch user:", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch user"})
		return
	}

	h.logger.Info("Fetched user successfully with ID:", id)
	c.JSON(http.StatusOK, u)
}

// DeleteUserById godoc
// @Summary      Delete a user
// @Description  Deletes a user by ID
// @Tags         Users
// @Param        id  path  int  true  "User ID"
// @Success      204
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /users/{id} [delete]
func (h *Handler) DeleteUserById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.logger.Error("Invalid user ID for deletion:", idStr)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	res, err := h.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		h.logger.Error("Failed to delete user:", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete user"})
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		h.logger.Warn("User not found for deletion with ID:", id)
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		return
	}

	h.logger.Info("Deleted user successfully with ID:", id)
	c.Status(http.StatusNoContent)
}