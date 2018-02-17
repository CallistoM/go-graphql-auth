package handler

import (
	// standard libraries
	"errors"
	// encryption library
	"golang.org/x/crypto/bcrypt"
	// custom handlers
	"github.com/callistom/go-graphql-auth/authentication"
	"github.com/callistom/go-graphql-auth/structs"
	// db
	"github.com/jinzhu/gorm"
)

// init vars
var (
	db  *gorm.DB
	err error
)

// Login function find user
func (r *Resolver) Login(args *struct {
	Input *structs.LoginInput
}) (string, error) {

	db, err = gorm.Open("postgres", "postgres://localhost/graphql?sslmode=disable")

	var users []structs.User

	if err := db.Find(&users).Error; err != nil {
		return "", errors.New("An error has occured when trying to fetch users")
	}

	for _, user := range users {
		if user.Mail == args.Input.Mail {

			// compare password with given password
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(args.Input.Password))

			if err == nil {
				token, err := authentication.GenerateToken(user)
				if err != nil {
					return "", err
				}
				return token, err
			}

			return "", errors.New("Password is not correct")
		}
	}

	return "", errors.New("User not found")
}
