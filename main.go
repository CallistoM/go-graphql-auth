package main

import (
	// standard libraries
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	// custom handlers
	authentication "github.com/callistom/go-graphql-auth/authentication"
	handler "github.com/callistom/go-graphql-auth/handler"
	structs "github.com/callistom/go-graphql-auth/structs"
	// db
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	// graphQL
	"github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/relay"
)

// init vars
var (
	db     *gorm.DB
	schema *graphql.Schema
	err    error
)

func init() {

	// read schema file
	schemaFile, err := ioutil.ReadFile("schema.graphql")

	if err != nil {
		panic(err)
	}

	// parse schema file
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

	// open db
	db, err = gorm.Open("postgres", "postgres://localhost/graphql?sslmode=disable")

	// migrate user struct
	db.AutoMigrate(&structs.User{})

	// check error
	if err != nil {
		log.Fatal(err)
	}

	// close db last
	defer db.Close()

	// set endpoint + authentication handler
	http.Handle("/graphql", authentication.Auth(&relay.Handler{Schema: schema}))

	fmt.Println("Listening at http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
