package handler

import (
	"errors"

	"github.com/callistom/go-graphql-auth/authentication"
	"github.com/callistom/go-graphql-auth/structs"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func (r *Resolver) Login(args *struct {
	Input *structs.LoginInput
}) (string, error) {

	var users []structs.User

	db.Find(&users)

	for _, user := range users {
		if user.Mail == args.Input.Mail {
			if user.Password == args.Input.Password {
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
