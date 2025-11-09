package main

import (
	"fmt"

	app "github.com/pratik6266/go-full/internal"

	"github.com/gin-gonic/gin"
	docs "github.com/pratik6266/go-full/docs" // generated swagger docs
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
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

	// Initialize DB once
	app.InitDB()

	// Initialize Gin router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(app.PrometheusMiddleware())

	// Swagger UI endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	docs.SwaggerInfo.BasePath = "/api/v1"

	rw := r.Group("/api/v1")
	{
		//healthcheck endpoint
		h := &app.Handler{}
		rw.GET("/health", h.Healthcheck)
		rw.GET("/students", h.GetStudents)
		rw.POST("/students", h.CreateStudent)
		rw.GET("/students/:id", h.GetStudentByID)
		rw.PUT("/students/:id", h.UpdateStudent)
		rw.DELETE("/students/:id", h.DeleteStudent)
	}

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Run(":8080")
}
