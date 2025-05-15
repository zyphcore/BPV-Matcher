package router

import (
	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	appConfig := config.LoadConfig()

	if appConfig.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{appConfig.AllowedOrigins}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization")
	router.Use(cors.New(corsConfig))

	jwtConfig := middleware.JWTConfig{
		SecretKey:     appConfig.JWTSecret,
		TokenDuration: appConfig.JWTExpiration,
	}

	// Initialize handlers
	userHandler := handlers.NewUserHandler(jwtConfig)
	adminHandler := handlers.NewAdminHandler()
	studentHandler := handlers.NewStudentHandler()
	coordinatorHandler := handlers.NewCoordinatorHandler()
	mentorHandler := handlers.NewMentorHandler()

	// Set up session timeout middleware
	sessionTimeout := middleware.SessionTimeoutMiddleware(8 * time.Hour)

	router.GET("/health", handlers.HealthCheck)

	// Authentication routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
	}

	// Common user routes
	user := router.Group("/api/user")
	user.Use(middleware.AuthMiddleware(jwtConfig), sessionTimeout)
	{
		user.GET("/profile", userHandler.GetProfile)
		user.PUT("/profile", userHandler.UpdateProfile)
		user.POST("/change-password", userHandler.ChangePassword)
	}

	// Student-specific routes
	student := router.Group("/api/student")
	student.Use(middleware.AuthMiddleware(jwtConfig), middleware.StudentMiddleware(), sessionTimeout)
	{
		student.GET("/profile", studentHandler.GetStudentProfile)
		student.PUT("/profile", studentHandler.UpdateStudentProfile)
		student.POST("/cv", studentHandler.UploadCV)
		student.POST("/cover-letter", studentHandler.UploadCoverLetter)
		student.PUT("/status", studentHandler.UpdateStatus)
	}

	// Coordinator-specific routes
	coordinator := router.Group("/api/coordinator")
	coordinator.Use(middleware.AuthMiddleware(jwtConfig), middleware.CoordinatorMiddleware(), sessionTimeout)
	{
		coordinator.GET("/profile", coordinatorHandler.GetCoordinatorProfile)
		coordinator.GET("/companies", coordinatorHandler.GetCompanies)
		coordinator.GET("/companies/:id", coordinatorHandler.GetCompany)
		coordinator.POST("/companies", coordinatorHandler.CreateCompany)
		coordinator.PUT("/companies/:id", coordinatorHandler.UpdateCompany)
		coordinator.PUT("/companies/:id/partner", coordinatorHandler.UpdateCompanyPartnerStatus)
		coordinator.GET("/statistics", coordinatorHandler.GetStudentStatistics)
		coordinator.GET("/students", coordinatorHandler.GetAllStudents)
	}

	// Mentor-specific routes
	mentor := router.Group("/api/mentor")
	mentor.Use(middleware.AuthMiddleware(jwtConfig), middleware.MentorMiddleware(), sessionTimeout)
	{
		mentor.GET("/profile", mentorHandler.GetMentorProfile)
		mentor.GET("/students/applications", mentorHandler.GetStudentApplications)
		mentor.GET("/students/inactive", mentorHandler.GetInactiveStudents)
		mentor.GET("/students/:id", mentorHandler.GetStudentDetails)
		mentor.POST("/students/:id/notify", mentorHandler.SendNotificationToStudent)
		mentor.GET("/dashboard", mentorHandler.GetStudentActivityDashboard)
	}

	// Admin routes
	admin := router.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(jwtConfig), middleware.AdminMiddleware(), sessionTimeout)
	{
		admin.GET("/users", adminHandler.GetAllUsers)
		admin.PUT("/users/:id/role", adminHandler.UpdateUserRole)
		admin.PUT("/users/:id/deactivate", adminHandler.DeactivateUser)
		admin.PUT("/users/:id/reactivate", adminHandler.ReactivateUser)
		admin.DELETE("/users/:id", adminHandler.DeleteUser)
		admin.POST("/users/:id/coordinator", adminHandler.AddCoordinatorPrivileges)
		admin.POST("/users/:id/mentor", adminHandler.AddMentorPrivileges)
		admin.DELETE("/users/:id/privileges", adminHandler.RemoveSpecialPrivileges)
	}

	return router
}
