package router

import (
	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewRouter creates and configures a new router
func NewRouter() *gin.Engine {
	// Load application configuration
	appConfig := config.LoadConfig()

	// Set Gin mode based on environment
	if appConfig.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	router := gin.Default()

	// Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{appConfig.AllowedOrigins}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization")
	router.Use(cors.New(corsConfig))

	// Create JWT configuration
	jwtConfig := middleware.JWTConfig{
		SecretKey:     appConfig.JWTSecret,
		TokenDuration: appConfig.JWTExpiration,
	}

	// Create handlers
	userHandler := handlers.NewUserHandler(jwtConfig)
	adminHandler := handlers.NewAdminHandler()

	// Public routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "BPV-Matcher API"})
	})

	// Auth routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
	}

	// User routes (protected)
	user := router.Group("/api/user")
	user.Use(middleware.AuthMiddleware(jwtConfig))
	{
		user.GET("/profile", userHandler.GetProfile)
		user.PUT("/profile", userHandler.UpdateProfile)
		user.POST("/change-password", userHandler.ChangePassword)
	}

	// Admin routes (protected + admin only)
	admin := router.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(jwtConfig), middleware.AdminMiddleware())
	{
		admin.GET("/users", adminHandler.GetAllUsers)
		admin.PUT("/users/:id/role", adminHandler.UpdateUserRole)
		admin.PUT("/users/:id/deactivate", adminHandler.DeactivateUser)
		admin.PUT("/users/:id/reactivate", adminHandler.ReactivateUser)
		admin.DELETE("/users/:id", adminHandler.DeleteUser)
		admin.POST("/users/:id/admin-privileges", adminHandler.AddAdminPrivileges)
		admin.DELETE("/users/:id/admin-privileges", adminHandler.RemoveAdminPrivileges)
	}

	return router
}
