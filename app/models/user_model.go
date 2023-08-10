package models

import (
	"time"

	"gorm.io/gorm"
)


type User struct {
	gorm.Model
	UserName string `gorm:"type:varchar(255)" json:"username"` 
	Password string `gorm:"type:varchar(255)" json:"password"`
	PhoneNumber string `gorm:"type:varchar(255)" json:"phonenumber"`
	IdFirebase string `gorm:"type:varchar(255)" json:"idfirebase"`


}

func (User) TableName() string {
	return "users"
}

type UserResponse struct {
	ID        uint      `json:"id,omitempty"`
	UserName      string    `json:"username" gorm:"type:varchar(255);not null"`
	PhoneNumber string `gorm:"type:varchar(255)" json:"phonenumber"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FilterUserRecord(user *User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		UserName:  user.UserName,
		PhoneNumber: user.PhoneNumber,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

type SignUpInput struct {
	PhoneNumber            string `json:"phonenumber" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
}

type SignInInput struct {
	PhoneNumber    string `json:"phonenumber"  validate:"required"`
	Password string `json:"password"  validate:"required"`
}
