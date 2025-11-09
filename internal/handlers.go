package internal

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

// Healthcheck godoc
// @Summary      Health check
// @Description  Returns the API health status.
// @Tags         Health
// @ID           healthcheck
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /health [get]
func (h *Handler) Healthcheck(c *gin.Context) {
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
	// Expect global DB to be initialized in main
	rows, err := DB.Query("SELECT id, name, age, email FROM students")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch students"})
		return
	}
	defer rows.Close()

	students := make([]Student, 0)
	for rows.Next() {
		var s Student
		if err := rows.Scan(&s.ID, &s.Name, &s.Age, &s.Email); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to parse student record"})
			return
		}
		students = append(students, s)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database error while fetching students"})
		return
	}

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
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	var student Student
	err := DB.QueryRow(
		"INSERT INTO students (name, age, email) VALUES ($1, $2, $3) RETURNING id",
		req.Name, req.Age, req.Email,
	).Scan(&student.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create student"})
		return
	}
	student.Name = req.Name
	student.Age = req.Age
	student.Email = req.Email

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
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	var s Student
	err = DB.QueryRow("SELECT id, name, age, email FROM students WHERE id = $1", id).
		Scan(&s.ID, &s.Name, &s.Age, &s.Email)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Student not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch student"})
		return
	}

	StudentCreationsTotal.WithLabelValues().Inc()
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
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	var payload StudentUpdateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	// Update and return the updated row
	var updated Student
	err = DB.QueryRow(
		"UPDATE students SET name = $1, age = $2, email = $3 WHERE id = $4 RETURNING id, name, age, email",
		payload.Name, payload.Age, payload.Email, id,
	).Scan(&updated.ID, &updated.Name, &updated.Age, &updated.Email)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Student not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update student"})
		return
	}

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
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	res, err := DB.Exec("DELETE FROM students WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete student"})
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Student not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
