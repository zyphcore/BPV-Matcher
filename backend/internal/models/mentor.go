package models

import (
	"gorm.io/gorm"
)

type Mentor struct {
	gorm.Model
	UserID                    uint `gorm:"uniqueIndex;not null"`
	User                      User `gorm:"foreignKey:UserID"`
	Department                string
	CanMonitorApplications    bool `gorm:"default:true"`
	CanTrackActivity          bool `gorm:"default:true"`
	CanNotifyInactiveStudents bool `gorm:"default:true"`
	AssignedStudentCount      int  `gorm:"default:0"`
}
