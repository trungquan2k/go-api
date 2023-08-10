package users

import (
	"gorm.io/gorm"
)

type PostApi struct {
	DB *gorm.DB
}

func InitPostApi(db *gorm.DB) *PostApi {
	return &PostApi{DB: db}
}
