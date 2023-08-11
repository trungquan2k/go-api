package models

import (
	"gorm.io/gorm"
)


type Ranks struct {
	gorm.Model
	Star string `gorm:"type:varchar(255)" json:"star"`
	AvgStar string `gorm:"type:varchar(255)" json:"avgStar"`
	Users []User
}
