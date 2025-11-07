package main

import (
	"fmt"

	handler "github.com/pratik6266/go-full/internal"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	docs "github.com/pratik6266/go-full/docs"  // generated swagger docs
)

// @title           Student API Documentation
// @version         1.0
// @description     This is a sample server for managing student records.

// @scheme http
// @baeshPath /api/v1

// @contact.name   Pratik Raj
// @contact.email pratikraj220011@gmail.com

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	fmt.Println("Server is running on port 8080")

	
	// Initialize Gin router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	
	// Swagger UI endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	docs.SwaggerInfo.BasePath = "/api/v1"

	// Initialize handlers
	h := &handler.Handler{}

	rw := r.Group("/api/v1")
	{
		rw.GET("/health", h.Healthcheck)
	}

	r.Run(":8080")
}
