package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	authentication "github.com/callistom/jwt_graphql_server/authentication"
	resolvers "github.com/callistom/jwt_graphql_server/resolvers"
	structs "github.com/callistom/jwt_graphql_server/structs"
	jwt "github.com/dgrijalva/jwt-go"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/relay"
)

type rResolver resolvers.Resolver

var db *gorm.DB

var schema *graphql.Schema

var users = []structs.User{
	{
		Name:     "Example User",
		Mail:     "example@mail.com",
		Password: "example",
	},
	{
		Name:     "Test User",
		Mail:     "test@mail.com",
		Password: "test",
	},
}

type User struct {
	gorm.Model
	Name     string
	Mail     string
	Password string
}

var err error

func init() {

	schemaFile, err := ioutil.ReadFile("schema.graphql")
	if err != nil {
		panic(err)
	}

	schema, err = graphql.ParseSchema(string(schemaFile), &rResolver{})
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

	db, err = gorm.Open("mysql", "root:admin@tcp(127.0.0.1)/graphql?charset=utf8&parseTime=True&loc=Local")

	db.AutoMigrate(&structs.User{})

	if err != nil {
		log.Fatal(err)
	}

	// defer db.Close()

	http.Handle("/graphql", authentication.Auth(&relay.Handler{Schema: schema}))

	fmt.Println("Listening at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (r *rResolver) Login(args *struct {
	Input *structs.LoginInput
}) (string, error) {
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

func (r *rResolver) Viewer(ctx context.Context, args *struct {
	Token *string
}) (*resolvers.UserResolver, error) {

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

	var user structs.User

	for _, u := range users {
		if id == string(u.ID) {
			user = u
		}
	}

	return &resolvers.UserResolver{
		User: user,
	}, nil
}

func (r *rResolver) Users(ctx context.Context, args *struct {
	Token *string
}) ([]*resolvers.UsersResolver, error) {
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

	// claims, _ := token.Claims.(jwt.MapClaims)
	// id := claims["sub"].(string)

	// var users = []structs.User{
	// 	{
	// 		ID:       "1",
	// 		Name:     "Example User",
	// 		Mail:     "example@mail.com",
	// 		Password: "example",
	// 	},
	// 	{
	// 		ID:       "2",
	// 		Name:     "Test User",
	// 		Mail:     "test@mail.com",
	// 		Password: "test",
	// 	},
	// }

	var allUsers []structs.User

	t := db.Find(&allUsers)

	for i, _ := range allUsers {
		db.Model(allUsers[i])
	}

	fmt.Println("places:", t)

	var l []*resolvers.UsersResolver

	for _, user := range users {
		l = append(l, &resolvers.UsersResolver{user})
	}

	return l, nil
}
