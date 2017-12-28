package structs

import "github.com/jinzhu/gorm"

// User struct
type User struct {
	gorm.Model
	Name     string
	Mail     string
	Password string
}

// LoginInput struct
type LoginInput struct {
	Mail     string
	Password string
}
