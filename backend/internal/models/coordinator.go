package models

import (
	"gorm.io/gorm"
)

type Coordinator struct {
	gorm.Model
	UserID             uint `gorm:"uniqueIndex;not null"`
	User               User `gorm:"foreignKey:UserID"`
	Department         string
	CanConfigSystem    bool `gorm:"default:true"`
	CanManageCompanies bool `gorm:"default:true"`
	CanViewStatistics  bool `gorm:"default:true"`
	CanMonitorStudents bool `gorm:"default:true"`
}
