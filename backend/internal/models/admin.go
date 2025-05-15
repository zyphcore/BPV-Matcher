package models

import (
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	UserID            uint   `gorm:"uniqueIndex;not null"`
	User              User   `gorm:"foreignKey:UserID"`
	AdminLevel        string `gorm:"default:standard;not null"`
	Department        string
	CanManageUsers    bool `gorm:"default:true"`
	CanApproveMatches bool `gorm:"default:true"`
	CanEditSettings   bool `gorm:"default:false"`
}
