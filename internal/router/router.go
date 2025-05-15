package router

import (
	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/middleware"

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

	userHandler := handlers.NewUserHandler(jwtConfig)
	adminHandler := handlers.NewAdminHandler()

	router.GET("/health", handlers.HealthCheck)

	auth := router.Group("/api/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
	}

	user := router.Group("/api/user")
	user.Use(middleware.AuthMiddleware(jwtConfig))
	{
		user.GET("/profile", userHandler.GetProfile)
		user.PUT("/profile", userHandler.UpdateProfile)
		user.POST("/change-password", userHandler.ChangePassword)
	}

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
