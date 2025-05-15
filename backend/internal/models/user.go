package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"uniqueIndex;not null"`
	Password     string `gorm:"not null"`
	FirstName    string `gorm:"not null"`
	LastName     string `gorm:"not null"`
	Role         string `gorm:"default:student;not null"`
	ProfileImage string
	Bio          string
	Active       bool `gorm:"default:true"`
	LastLogin    time.Time
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.Password == "" {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// Helper functions to check user roles
func (u *User) IsStudent() bool {
	return u.Role == "student"
}

func (u *User) IsCoordinator() bool {
	return u.Role == "coordinator"
}

func (u *User) IsMentor() bool {
	return u.Role == "mentor"
}
