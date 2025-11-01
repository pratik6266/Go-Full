package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	handler "github.com/pratik6266/go-full/internal"
)

func main() {
	fmt.Println("Server is running on port 8080")

	// Initialize Gin router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Initialize handlers
	h := &handler.Handler{}

	rw := r.Group("/api/v1")
	{
		rw.GET("/health", h.Healthcheck)
	}

	r.Run(":8080")
}
