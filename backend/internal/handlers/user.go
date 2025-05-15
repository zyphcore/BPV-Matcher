package handlers

import (
	"net/http"
	"strconv"
	"time"

	"backend/internal/database"
	"backend/internal/middleware"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	jwtConfig middleware.JWTConfig
}

func NewUserHandler(jwtConfig middleware.JWTConfig) *UserHandler {
	return &UserHandler{
		jwtConfig: jwtConfig,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var registrationData struct {
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=8"`
		FirstName string `json:"firstName" binding:"required"`
		LastName  string `json:"lastName" binding:"required"`
		Role      string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&registrationData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validRoles := map[string]bool{
		"student":     true,
		"coordinator": true,
		"mentor":      true,
	}

	if !validRoles[registrationData.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Must be student, coordinator, or mentor."})
		return
	}

	var existingUser models.User
	result := database.GetDB().Where("email = ?", registrationData.Email).First(&existingUser)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already in use"})
		return
	}

	user := models.User{
		Email:     registrationData.Email,
		Password:  registrationData.Password,
		FirstName: registrationData.FirstName,
		LastName:  registrationData.LastName,
		Role:      registrationData.Role,
		Active:    true,
		LastLogin: time.Now(),
	}

	result = database.GetDB().Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	switch user.Role {
	case "student":
		student := models.Student{
			UserID: user.ID,
		}
		database.GetDB().Create(&student)
	case "coordinator":
		coordinator := models.Coordinator{
			UserID: user.ID,
		}
		database.GetDB().Create(&coordinator)
	case "mentor":
		mentor := models.Mentor{
			UserID: user.ID,
		}
		database.GetDB().Create(&mentor)
	}

	user.Password = ""

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": user})
}

func (h *UserHandler) Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	result := database.GetDB().Where("email = ?", loginData.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !user.CheckPassword(loginData.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !user.Active {
		c.JSON(http.StatusForbidden, gin.H{"error": "Your account has been deactivated. Please contact an administrator."})
		return
	}

	token, err := middleware.GenerateToken(user, h.jwtConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	database.GetDB().Model(&user).Update("last_login", time.Now())

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"role":      user.Role,
		},
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	var user models.User
	result := database.GetDB().First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Password = ""

	switch role {
	case "student":
		var student models.Student
		database.GetDB().Where("user_id = ?", userID).First(&student)
		c.JSON(http.StatusOK, gin.H{"user": user, "studentProfile": student})
		return
	case "coordinator":
		var coordinator models.Coordinator
		database.GetDB().Where("user_id = ?", userID).First(&coordinator)
		c.JSON(http.StatusOK, gin.H{"user": user, "coordinatorProfile": coordinator})
		return
	case "mentor":
		var mentor models.Mentor
		database.GetDB().Where("user_id = ?", userID).First(&mentor)
		c.JSON(http.StatusOK, gin.H{"user": user, "mentorProfile": mentor})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("userID")

	var updateData struct {
		FirstName    string `json:"firstName"`
		LastName     string `json:"lastName"`
		ProfileImage string `json:"profileImage"`
		Bio          string `json:"bio"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	result := database.GetDB().First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if updateData.FirstName != "" {
		user.FirstName = updateData.FirstName
	}
	if updateData.LastName != "" {
		user.LastName = updateData.LastName
	}
	if updateData.ProfileImage != "" {
		user.ProfileImage = updateData.ProfileImage
	}
	if updateData.Bio != "" {
		user.Bio = updateData.Bio
	}

	result = database.GetDB().Save(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully", "user": user})
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("userID")

	var passwordData struct {
		CurrentPassword string `json:"currentPassword" binding:"required"`
		NewPassword     string `json:"newPassword" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&passwordData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	result := database.GetDB().First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if !user.CheckPassword(passwordData.CurrentPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
		return
	}

	user.Password = passwordData.NewPassword
	result = database.GetDB().Save(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	result := database.GetDB().First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{"user": user})
}
