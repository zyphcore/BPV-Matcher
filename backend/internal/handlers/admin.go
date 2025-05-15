package handlers

import (
	"net/http"
	"strconv"

	"backend/internal/database"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) GetAllUsers(c *gin.Context) {
	var users []models.User
	result := database.GetDB().Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var updateData struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updateData.Role != "user" && updateData.Role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	result := database.GetDB().Model(&models.User{}).Where("id = ?", userID).Update("role", updateData.Role)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User role updated successfully"})
}

func (h *AdminHandler) DeactivateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	result := database.GetDB().Model(&models.User{}).Where("id = ?", userID).Update("active", false)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate user"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deactivated successfully"})
}

func (h *AdminHandler) ReactivateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	result := database.GetDB().Model(&models.User{}).Where("id = ?", userID).Update("active", true)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reactivate user"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User reactivated successfully"})
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	result := database.GetDB().Delete(&models.User{}, userID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *AdminHandler) AddAdminPrivileges(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	result := database.GetDB().First(&user, userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		}
		return
	}

	var admin models.Admin
	result = database.GetDB().Where("user_id = ?", userID).First(&admin)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already has admin privileges"})
		return
	}

	admin = models.Admin{
		UserID:            uint(userID),
		AdminLevel:        "standard",
		CanManageUsers:    true,
		CanApproveMatches: true,
		CanEditSettings:   false,
	}

	result = database.GetDB().Create(&admin)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add admin privileges"})
		return
	}

	database.GetDB().Model(&user).Update("role", "admin")

	c.JSON(http.StatusOK, gin.H{"message": "Admin privileges added successfully"})
}

func (h *AdminHandler) RemoveAdminPrivileges(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	result := database.GetDB().Where("user_id = ?", userID).Delete(&models.Admin{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove admin privileges"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin record not found"})
		return
	}

	database.GetDB().Model(&models.User{}).Where("id = ?", userID).Update("role", "user")

	c.JSON(http.StatusOK, gin.H{"message": "Admin privileges removed successfully"})
}
