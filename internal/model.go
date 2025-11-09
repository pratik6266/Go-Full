package internal

type HealthResponse struct {
	Message string `json:"message" example:"API is healthy"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"internal server error"`
}

type Student struct {
	ID    int    `json:"id" example:"1"`
	Name  string `json:"name" example:"John Doe"`
	Age   int    `json:"age" example:"20"`
	Email string `json:"email" example:"john@example.com"`
}

// StudentCreateRequest represents the payload to create a new student.
type StudentCreateRequest struct {
	Name  string `json:"name" example:"John Doe"`
	Age   int    `json:"age" example:"20"`
	Email string `json:"email" example:"john@example.com"`
}

// StudentUpdateRequest represents the payload to update an existing student.
type StudentUpdateRequest struct {
	Name  string `json:"name" example:"John Doe"`
	Age   int    `json:"age" example:"21"`
	Email string `json:"email" example:"john.new@example.com"`
}

type User struct {
	ID    int    `json:"id" example:"1"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
}

type UserCreateRequest struct {
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
}
