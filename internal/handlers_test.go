package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// helper to create a handler with a mocked DB
func newTestHandler(t *testing.T) (*Handler, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	d := &Db{db: mockDB}
	logger := logrus.New()
	return NewHandler(d, logger), mock
}

func TestHealthcheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &Handler{db: nil, logger: logrus.New()}
	r := gin.New()
	r.GET("/health", h.Healthcheck)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp HealthResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "API is healthy", resp.Message)
}

func TestGetUsersSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, mock := newTestHandler(t)
	r := gin.New()
	r.GET("/users", h.GetUsers)

	rows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(1, "Alice", "alice@example.com").
		AddRow(2, "Bob", "bob@example.com")
	mock.ExpectQuery("SELECT id, name, email FROM users").WillReturnRows(rows)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var users []User
	_ = json.Unmarshal(w.Body.Bytes(), &users)
	assert.Len(t, users, 2)
	assert.Equal(t, "Alice", users[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByIdInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, _ := newTestHandler(t)
	r := gin.New()
	r.GET("/users/by-id", h.GetUserById)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users/by-id?id=abc", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUserByIdNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, mock := newTestHandler(t)
	r := gin.New()
	r.GET("/users/by-id", h.GetUserById)

	rows := sqlmock.NewRows([]string{"id", "name", "email"}) // empty result set
	mock.ExpectQuery("SELECT id, name, email FROM users WHERE id = \\$1").WithArgs(99).WillReturnRows(rows)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users/by-id?id=99", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUserByIdSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, mock := newTestHandler(t)
	r := gin.New()
	r.DELETE("/users/:id", h.DeleteUserById)

	mock.ExpectExec("DELETE FROM users WHERE id = \\$1").WithArgs(7).WillReturnResult(sqlmock.NewResult(0, 1))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/users/7", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUserByIdNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, mock := newTestHandler(t)
	r := gin.New()
	r.DELETE("/users/:id", h.DeleteUserById)

	mock.ExpectExec("DELETE FROM users WHERE id = \\$1").WithArgs(8).WillReturnResult(sqlmock.NewResult(0, 0))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/users/8", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}
