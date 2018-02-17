package structs

import "time"

// User struct
type User struct {
	ID        uint `gorm:"primary_key" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Name      string
	Mail      string
	Password  string
}

// LoginInput struct
type LoginInput struct {
	Mail     string
	Password string
}
