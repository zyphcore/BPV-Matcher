package models

import (
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	UserID               uint `gorm:"uniqueIndex;not null"`
	User                 User `gorm:"foreignKey:UserID"`
	ProgrammingLanguages string
	Preferences          string
	PortfolioURL         string
	CV                   string // Path to CV file
	CoverLetter          string // Path to cover letter file
	Status               string `gorm:"default:inactive"` // Status: inactive, searching, applied, placed
}
