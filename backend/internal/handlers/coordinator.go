package handlers

import (
	"net/http"
	"strconv"

	"backend/internal/database"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CoordinatorHandler struct{}

func NewCoordinatorHandler() *CoordinatorHandler {
	return &CoordinatorHandler{}
}

func (h *CoordinatorHandler) GetCoordinatorProfile(c *gin.Context) {
	userID, _ := c.Get("userID")

	var coordinator models.Coordinator
	result := database.GetDB().Where("user_id = ?", userID).Preload("User").First(&coordinator)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coordinator profile not found"})
		return
	}

	coordinator.User.Password = ""

	c.JSON(http.StatusOK, gin.H{"coordinator": coordinator})
}

func (h *CoordinatorHandler) GetCompanies(c *gin.Context) {
	var companies []models.Company
	result := database.GetDB().Find(&companies)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve companies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"companies": companies})
}

func (h *CoordinatorHandler) GetCompany(c *gin.Context) {
	companyIDStr := c.Param("id")
	companyID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	var company models.Company
	result := database.GetDB().First(&company, companyID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve company"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"company": company})
}

func (h *CoordinatorHandler) CreateCompany(c *gin.Context) {
	var company models.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.GetDB().Create(&company)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create company"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Company created successfully", "company": company})
}

func (h *CoordinatorHandler) UpdateCompany(c *gin.Context) {
	companyIDStr := c.Param("id")
	companyID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	var company models.Company
	result := database.GetDB().First(&company, companyID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve company"})
		}
		return
	}

	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result = database.GetDB().Save(&company)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update company"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Company updated successfully", "company": company})
}

func (h *CoordinatorHandler) UpdateCompanyPartnerStatus(c *gin.Context) {
	companyIDStr := c.Param("id")
	companyID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	var updateData struct {
		IsPartner bool `json:"isPartner"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.GetDB().Model(&models.Company{}).Where("id = ?", companyID).Update("is_partner", updateData.IsPartner)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update company partner status"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	status := "partner"
	if !updateData.IsPartner {
		status = "non-partner"
	}

	c.JSON(http.StatusOK, gin.H{"message": "Company updated to " + status + " successfully"})
}

func (h *CoordinatorHandler) GetStudentStatistics(c *gin.Context) {
	var stats struct {
		Total     int64 `json:"total"`
		Inactive  int64 `json:"inactive"`
		Searching int64 `json:"searching"`
		Applied   int64 `json:"applied"`
		Placed    int64 `json:"placed"`
	}

	database.GetDB().Model(&models.Student{}).Count(&stats.Total)
	database.GetDB().Model(&models.Student{}).Where("status = ?", "inactive").Count(&stats.Inactive)
	database.GetDB().Model(&models.Student{}).Where("status = ?", "searching").Count(&stats.Searching)
	database.GetDB().Model(&models.Student{}).Where("status = ?", "applied").Count(&stats.Applied)
	database.GetDB().Model(&models.Student{}).Where("status = ?", "placed").Count(&stats.Placed)

	c.JSON(http.StatusOK, gin.H{"statistics": stats})
}

func (h *CoordinatorHandler) GetAllStudents(c *gin.Context) {
	var students []models.Student
	result := database.GetDB().Preload("User").Find(&students)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve students"})
		return
	}

	for i := range students {
		students[i].User.Password = ""
	}

	c.JSON(http.StatusOK, gin.H{"students": students})
}
