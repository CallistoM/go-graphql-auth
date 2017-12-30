package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	authentication "github.com/callistom/jwt_graphql_server/authentication"
	handler "github.com/callistom/jwt_graphql_server/handler"
	structs "github.com/callistom/jwt_graphql_server/structs"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/relay"
)

var db *gorm.DB

var schema *graphql.Schema

var err error

func init() {

	schemaFile, err := ioutil.ReadFile("schema.graphql")
	if err != nil {
		panic(err)
	}

	schema, err = graphql.ParseSchema(string(schemaFile), &handler.Resolver{})
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

	defer db.Close()

	http.Handle("/graphql", authentication.Auth(&relay.Handler{Schema: schema}))

	fmt.Println("Listening at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// func (r *rResolver) User(ctx context.Context, args *struct {
// 	Token *string
// }) (*resolvers.UserResolver, error) {
// 	token := ctx.Value("jwt").(*jwt.Token)

// 	if token == nil && args.Token == nil {
// 		return nil, errors.New("Token not set")
// 	}

// 	if token == nil && args.Token != nil {
// 		viewerToken, err := authentication.CheckToken(*args.Token)
// 		if err != nil {
// 			return nil, err
// 		}
// 		token = viewerToken
// 	}

// 	claims := token.Claims.(*authentication.MyCustomClaims)

// 	var user structs.User

// 	var users []structs.User

// 	db.Find(&users)

// 	for _, u := range users {
// 		if claims.ID == string(u.ID) {
// 			user = u
// 		}
// 	}

// 	return &resolvers.UserResolver{
// 		User: user,
// 	}, nil
// }

// func (r *rResolver) Users(ctx context.Context, args *struct {
// 	Token *string
// }) ([]*resolvers.UsersResolver, error) {
// 	token := ctx.Value("jwt").(*jwt.Token)

// 	if token == nil && args.Token == nil {
// 		return nil, errors.New("There needs to be a token in the Authorization header or viewer input")
// 	}

// 	if token == nil && args.Token != nil {
// 		viewerToken, err := authentication.CheckToken(*args.Token)
// 		if err != nil {
// 			return nil, err
// 		}
// 		token = viewerToken
// 	}

// 	var allUsers []structs.User

// 	db.Find(&allUsers)

// 	var users []*resolvers.UsersResolver

// 	for _, u := range allUsers {
// 		users = append(users, &resolvers.UsersResolver{u})
// 	}

// 	return users, nil
// }
