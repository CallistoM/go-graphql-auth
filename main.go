package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"io/ioutil"

	"fmt"

	authentication "github.com/callistom/api-project/authentication"
	userStruct "github.com/callistom/api-project/structs"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/relay"
)

type LoginInput struct {
	Mail     string
	Password string
}

var users = []userStruct.User{
	{
		ID:       "1",
		Name:     "Example User",
		Mail:     "example@mail.com",
		Password: "example",
	},
	{
		ID:       "2",
		Name:     "Test User",
		Mail:     "test@mail.com",
		Password: "test",
	},
}

type viewerResolver struct {
	User userStruct.User
}

var schema *graphql.Schema

func init() {

	schemaFile, err := ioutil.ReadFile("schema.graphql")
	if err != nil {
		panic(err)
	}

	schema, err = graphql.ParseSchema(string(schemaFile), &Resolver{})
	if err != nil {
		panic(err)
	}
}

func main() {
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, err := ioutil.ReadFile("graphiql.html")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(page)
	}))

	http.Handle("/graphql", authentication.Auth(&relay.Handler{Schema: schema}))

	fmt.Println("Listening at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Resolver struct{}

func (r *Resolver) Login(args *struct {
	Input *LoginInput
}) (string, error) {
	for _, user := range users {
		if user.Mail == args.Input.Mail {
			if user.Password == args.Input.Password {
				token, err := authentication.GenerateToken(user)
				if err != nil {
					return "", err
				}
				return token, err
			} else {
				return "", errors.New("password is incorrect")
			}
		}
	}

	return "", errors.New("User not found")
}

func (r *Resolver) Viewer(ctx context.Context, args *struct {
	Token *string
}) (*viewerResolver, error) {

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

	claims, _ := token.Claims.(jwt.MapClaims)
	id := claims["sub"].(string)

	var user userStruct.User

	for _, u := range users {
		if id == string(u.ID) {
			user = u
		}
	}

	return &viewerResolver{
		User: user,
	}, nil
}

func (v *viewerResolver) ID() graphql.ID {
	return graphql.ID(v.User.ID)
}

func (v *viewerResolver) Name() string {
	return v.User.Name
}

func (v *viewerResolver) Mail() string {
	return v.User.Mail
}
