package internal

type HealthResponse struct {
	Message string `json:"message" example:"API is healthy"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"internal server error"`
}
