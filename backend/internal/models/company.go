package models

import (
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	Website     string
	Location    string
	Industry    string
	IsPartner   bool `gorm:"default:false"`
	IsPotential bool `gorm:"default:false"`
	Contact     string
	Email       string
	Phone       string
	Notes       string
}
