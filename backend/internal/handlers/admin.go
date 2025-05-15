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

	// Validate the role
	validRoles := map[string]bool{
		"student":     true,
		"coordinator": true,
		"mentor":      true,
	}

	if !validRoles[updateData.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Must be student, coordinator, or mentor."})
		return
	}

	var user models.User
	result := database.GetDB().First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Store the old role for comparison
	oldRole := user.Role

	// Update the user's role
	result = database.GetDB().Model(&user).Update("role", updateData.Role)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
		return
	}

	// Create role-specific profile if it doesn't exist
	if oldRole != updateData.Role {
		switch updateData.Role {
		case "student":
			var student models.Student
			if database.GetDB().Where("user_id = ?", userID).First(&student).Error != nil {
				student = models.Student{UserID: uint(userID)}
				database.GetDB().Create(&student)
			}
		case "coordinator":
			var coordinator models.Coordinator
			if database.GetDB().Where("user_id = ?", userID).First(&coordinator).Error != nil {
				coordinator = models.Coordinator{UserID: uint(userID)}
				database.GetDB().Create(&coordinator)
			}
		case "mentor":
			var mentor models.Mentor
			if database.GetDB().Where("user_id = ?", userID).First(&mentor).Error != nil {
				mentor = models.Mentor{UserID: uint(userID)}
				database.GetDB().Create(&mentor)
			}
		}
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

	// Begin a transaction to ensure all related records are deleted
	tx := database.GetDB().Begin()

	// Delete role-specific profile if it exists
	var user models.User
	if tx.First(&user, userID).Error == nil {
		switch user.Role {
		case "student":
			tx.Where("user_id = ?", userID).Delete(&models.Student{})
		case "coordinator":
			tx.Where("user_id = ?", userID).Delete(&models.Coordinator{})
		case "mentor":
			tx.Where("user_id = ?", userID).Delete(&models.Mentor{})
		}
	}

	// Delete the user
	result := tx.Delete(&models.User{}, userID)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *AdminHandler) AddCoordinatorPrivileges(c *gin.Context) {
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

	var coordinator models.Coordinator
	result = database.GetDB().Where("user_id = ?", userID).First(&coordinator)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already has coordinator privileges"})
		return
	}

	coordinator = models.Coordinator{
		UserID:             uint(userID),
		Department:         "General",
		CanConfigSystem:    true,
		CanManageCompanies: true,
		CanViewStatistics:  true,
		CanMonitorStudents: true,
	}

	result = database.GetDB().Create(&coordinator)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add coordinator privileges"})
		return
	}

	database.GetDB().Model(&user).Update("role", "coordinator")

	c.JSON(http.StatusOK, gin.H{"message": "Coordinator privileges added successfully"})
}

func (h *AdminHandler) AddMentorPrivileges(c *gin.Context) {
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

	var mentor models.Mentor
	result = database.GetDB().Where("user_id = ?", userID).First(&mentor)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already has mentor privileges"})
		return
	}

	mentor = models.Mentor{
		UserID:                    uint(userID),
		Department:                "General",
		CanMonitorApplications:    true,
		CanTrackActivity:          true,
		CanNotifyInactiveStudents: true,
	}

	result = database.GetDB().Create(&mentor)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add mentor privileges"})
		return
	}

	database.GetDB().Model(&user).Update("role", "mentor")

	c.JSON(http.StatusOK, gin.H{"message": "Mentor privileges added successfully"})
}

func (h *AdminHandler) RemoveSpecialPrivileges(c *gin.Context) {
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

	switch user.Role {
	case "coordinator":
		database.GetDB().Where("user_id = ?", userID).Delete(&models.Coordinator{})
	case "mentor":
		database.GetDB().Where("user_id = ?", userID).Delete(&models.Mentor{})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "User does not have special privileges"})
		return
	}

	database.GetDB().Model(&user).Update("role", "student")

	// Create student profile if it doesn't exist
	var student models.Student
	if database.GetDB().Where("user_id = ?", userID).First(&student).Error != nil {
		student = models.Student{UserID: uint(userID)}
		database.GetDB().Create(&student)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Special privileges removed successfully"})
}
