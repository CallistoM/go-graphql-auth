package handler

import (
	// standard library
	"context"
	"errors"
	// custom handlers
	"github.com/callistom/go-graphql-auth/authentication"
	"github.com/callistom/go-graphql-auth/resolver"
	"github.com/callistom/go-graphql-auth/structs"
	"github.com/jinzhu/gorm"
	// jwt package
	jwt "github.com/dgrijalva/jwt-go"
)

// User handler
func (r *Resolver) User(ctx context.Context, args *struct {
	Token *string
}) (*resolver.UserResolver, error) {

	db, err = gorm.Open("postgres", "postgres://localhost/graphql?sslmode=disable")

	token := ctx.Value("jwt").(*jwt.Token)

	if token == nil && args.Token == nil {
		return nil, errors.New("Token not set")
	}

	if token == nil && args.Token != nil {
		viewerToken, err := authentication.CheckToken(*args.Token)

		if err != nil {
			return nil, err
		}

		token = viewerToken
	}

	claims := token.Claims.(*authentication.MyCustomClaims)

	var (
		user  structs.User
		users []structs.User
	)

	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	for _, u := range users {
		if claims.ID == string(u.ID) {
			user = u
		}
	}

	return &resolver.UserResolver{
		User: user,
	}, nil
}

// Users handler
func (r *Resolver) Users(ctx context.Context, args *struct {
	Token *string
}) ([]*resolver.UsersResolver, error) {

	token := ctx.Value("jwt").(*jwt.Token)

	if token == nil && args.Token == nil {
		return nil, errors.New("There needs to be a token in the Authorization header or viewer input")
	}

	if token == nil && args.Token != nil {
		viewerToken, err := authentication.CheckToken(*args.Token)
		if err != nil {
			return nil, err
		}
		token = viewerToken
	}

	var (
		allUsers []structs.User
		users    []*resolver.UsersResolver
	)

	db.Find(&allUsers)

	for _, u := range allUsers {
		users = append(users, &resolver.UsersResolver{u})
	}

	return users, nil
}
