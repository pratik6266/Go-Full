package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grafana/loki-client-go/loki"
	docs "github.com/pratik6266/go-full/docs" // generated swagger docs
	app "github.com/pratik6266/go-full/internal"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// LokiHook implements logrus.Hook to forward logs to Loki.
type LokiHook struct {
	client *loki.Client
	labels model.LabelSet
}

func NewLokiHook(client *loki.Client) *LokiHook {
	return &LokiHook{
		client: client,
		labels: model.LabelSet{model.LabelName("job"): model.LabelValue("student-api")},
	}
}

func (h *LokiHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *LokiHook) Fire(entry *logrus.Entry) error {
	// Format log line
	line, err := entry.String()
	if err != nil {
		return err
	}
	// Send to Loki; ignore error for now to avoid loop
	_ = h.client.Handle(h.labels, time.Now(), line)
	return nil
}

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

	// Initialize database
	db, err := app.InitDB()
	if err != nil {
		panic("Failed to connect to database")
	}

	// Initialize Gin router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(app.PrometheusMiddleware())

	// Loki Setup
	lokiCfg, err := loki.NewDefaultConfig("http://localhost:3100/loki/api/v1/push")
	if err != nil {
		panic(err)
	}
	client, err := loki.New(lokiCfg)
	if err != nil {
		panic(err)
	}
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})
	logger.AddHook(NewLokiHook(client))

	// Swagger UI endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	docs.SwaggerInfo.BasePath = "/api/v1"

	// API routes
	h := app.NewHandler(db, logger)
	rw := r.Group("/api/v1")
	{
		//healthcheck endpoint
		rw.GET("/health", h.Healthcheck)

		// student endpoints
		rw.GET("/students", h.GetStudents)
		rw.POST("/students", h.CreateStudent)
		rw.GET("/students/:id", h.GetStudentByID)
		rw.PUT("/students/:id", h.UpdateStudent)
		rw.DELETE("/students/:id", h.DeleteStudent)

		// user endpoints
		rw.GET("/users", h.GetUsers)
		rw.GET("/users/by-id", h.GetUserById)
		rw.POST("/users", h.CreateUser)
		rw.DELETE("/users/:id", h.DeleteUserById)

		// auth endpoints
	}

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Run(":8080")
}
