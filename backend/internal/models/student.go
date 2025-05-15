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
	CV                   string
	CoverLetter          string
	Status               string `gorm:"default:inactive"`
}
