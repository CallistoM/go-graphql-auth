package handler

import (
	"context"
	"errors"

	"github.com/callistom/go-graphql-auth/authentication"
	"github.com/callistom/go-graphql-auth/resolver"
	"github.com/callistom/go-graphql-auth/structs"
	jwt "github.com/dgrijalva/jwt-go"
)

func (r *Resolver) User(ctx context.Context, args *struct {
	Token *string
}) (*resolver.UserResolver, error) {
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

	var user structs.User

	var users []structs.User

	db.Find(&users)

	for _, u := range users {
		if claims.ID == string(u.ID) {
			user = u
		}
	}

	return &resolver.UserResolver{
		User: user,
	}, nil
}
