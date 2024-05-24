package dbmodel

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"not null"`
	Hash     string `gorm:"not null" json:"-"`
	Role     string `gorm:"not null" json:"-"`
}
