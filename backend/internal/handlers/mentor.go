package handlers

import (
	"net/http"
	"strconv"
	"time"

	"backend/internal/database"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
)

type MentorHandler struct{}

func NewMentorHandler() *MentorHandler {
	return &MentorHandler{}
}

// GetMentorProfile retrieves the mentor's profile
func (h *MentorHandler) GetMentorProfile(c *gin.Context) {
	userID, _ := c.Get("userID")

	var mentor models.Mentor
	result := database.GetDB().Where("user_id = ?", userID).Preload("User").First(&mentor)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mentor profile not found"})
		return
	}

	mentor.User.Password = ""

	c.JSON(http.StatusOK, gin.H{"mentor": mentor})
}

// GetStudentApplications retrieves student applications
func (h *MentorHandler) GetStudentApplications(c *gin.Context) {
	var students []models.Student
	result := database.GetDB().Where("status = ? OR status = ?", "applied", "searching").Preload("User").Find(&students)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve students"})
		return
	}

	// Remove sensitive information
	for i := range students {
		students[i].User.Password = ""
	}

	c.JSON(http.StatusOK, gin.H{"students": students})
}

// GetInactiveStudents retrieves inactive students
func (h *MentorHandler) GetInactiveStudents(c *gin.Context) {
	var students []models.Student
	oneMonthAgo := time.Now().AddDate(0, -1, 0)

	result := database.GetDB().
		Joins("JOIN users ON students.user_id = users.id").
		Where("students.status = ? OR (students.status = ? AND users.last_login < ?)",
			"inactive", "searching", oneMonthAgo).
		Preload("User").
		Find(&students)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve inactive students"})
		return
	}

	// Remove sensitive information
	for i := range students {
		students[i].User.Password = ""
	}

	c.JSON(http.StatusOK, gin.H{"students": students})
}

// GetStudentDetails retrieves a specific student's details
func (h *MentorHandler) GetStudentDetails(c *gin.Context) {
	studentIDStr := c.Param("id")
	studentID, err := strconv.ParseUint(studentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var student models.Student
	result := database.GetDB().Preload("User").First(&student, studentID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	student.User.Password = ""

	c.JSON(http.StatusOK, gin.H{"student": student})
}

// SendNotificationToStudent sends a notification to an inactive student
func (h *MentorHandler) SendNotificationToStudent(c *gin.Context) {
	studentIDStr := c.Param("id")
	studentID, err := strconv.ParseUint(studentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var notificationData struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&notificationData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real application, you would send an email or push notification
	// For this example, we'll just return a success message

	// You could also store the notification in a database

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification sent successfully",
		"details": gin.H{
			"studentID":    studentID,
			"notification": notificationData.Message,
		},
	})
}

// GetStudentActivityDashboard gets activity data for dashboard
func (h *MentorHandler) GetStudentActivityDashboard(c *gin.Context) {
	// Count students by status for the dashboard
	var stats struct {
		Total          int64 `json:"total"`
		Inactive       int64 `json:"inactive"`
		Searching      int64 `json:"searching"`
		Applied        int64 `json:"applied"`
		Placed         int64 `json:"placed"`
		RecentlyActive int64 `json:"recentlyActive"`
	}

	database.GetDB().Model(&models.Student{}).Count(&stats.Total)
	database.GetDB().Model(&models.Student{}).Where("status = ?", "inactive").Count(&stats.Inactive)
	database.GetDB().Model(&models.Student{}).Where("status = ?", "searching").Count(&stats.Searching)
	database.GetDB().Model(&models.Student{}).Where("status = ?", "applied").Count(&stats.Applied)
	database.GetDB().Model(&models.Student{}).Where("status = ?", "placed").Count(&stats.Placed)

	// Count recently active students (logged in within the last week)
	oneWeekAgo := time.Now().AddDate(0, 0, -7)
	database.GetDB().
		Model(&models.User{}).
		Joins("JOIN students ON users.id = students.user_id").
		Where("users.last_login > ?", oneWeekAgo).
		Count(&stats.RecentlyActive)

	c.JSON(http.StatusOK, gin.H{"statistics": stats})
}
