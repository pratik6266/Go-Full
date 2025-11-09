package internal

import(
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)


func HealthcheckTest(t *testing.T) {
	router := gin.Default()

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		panic(err);
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}