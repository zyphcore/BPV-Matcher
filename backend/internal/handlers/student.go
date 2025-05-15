package handlers

import (
	"net/http"
	"strconv"

	"backend/internal/database"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
)

type StudentHandler struct{}

func NewStudentHandler() *StudentHandler {
	return &StudentHandler{}
}

func (h *StudentHandler) GetStudentProfile(c *gin.Context) {
	userID, _ := c.Get("userID")

	var student models.Student
	result := database.GetDB().Where("user_id = ?", userID).Preload("User").First(&student)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student profile not found"})
		return
	}

	student.User.Password = ""

	c.JSON(http.StatusOK, gin.H{"student": student})
}

func (h *StudentHandler) UpdateStudentProfile(c *gin.Context) {
	userID, _ := c.Get("userID")

	var updateData struct {
		ProgrammingLanguages string `json:"programmingLanguages"`
		Preferences          string `json:"preferences"`
		PortfolioURL         string `json:"portfolioURL"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var student models.Student
	result := database.GetDB().Where("user_id = ?", userID).First(&student)

	if result.Error != nil {
		student = models.Student{
			UserID:               uint(userID.(uint)),
			ProgrammingLanguages: updateData.ProgrammingLanguages,
			Preferences:          updateData.Preferences,
			PortfolioURL:         updateData.PortfolioURL,
		}
		result = database.GetDB().Create(&student)
	} else {
		if updateData.ProgrammingLanguages != "" {
			student.ProgrammingLanguages = updateData.ProgrammingLanguages
		}
		if updateData.Preferences != "" {
			student.Preferences = updateData.Preferences
		}
		if updateData.PortfolioURL != "" {
			student.PortfolioURL = updateData.PortfolioURL
		}
		result = database.GetDB().Save(&student)
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student profile updated successfully", "student": student})
}

func (h *StudentHandler) UploadCV(c *gin.Context) {
	userID, _ := c.Get("userID")

	file, err := c.FormFile("cv")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	filePath := "uploads/cv/user_" + strconv.FormatUint(uint64(userID.(uint)), 10) + "_" + file.Filename

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	result := database.GetDB().Model(&models.Student{}).Where("user_id = ?", userID).Update("cv", filePath)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update CV file path"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "CV uploaded successfully", "filePath": filePath})
}

func (h *StudentHandler) UploadCoverLetter(c *gin.Context) {
	userID, _ := c.Get("userID")

	file, err := c.FormFile("coverLetter")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	filePath := "uploads/cover_letters/user_" + strconv.FormatUint(uint64(userID.(uint)), 10) + "_" + file.Filename

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	result := database.GetDB().Model(&models.Student{}).Where("user_id = ?", userID).Update("cover_letter", filePath)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cover letter file path"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cover letter uploaded successfully", "filePath": filePath})
}

func (h *StudentHandler) UpdateStatus(c *gin.Context) {
	userID, _ := c.Get("userID")

	var updateData struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validStatuses := map[string]bool{
		"inactive":  true,
		"searching": true,
		"applied":   true,
		"placed":    true,
	}

	if !validStatuses[updateData.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	result := database.GetDB().Model(&models.Student{}).Where("user_id = ?", userID).Update("status", updateData.Status)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}
