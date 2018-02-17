package main

import (
	// standard libraries
	"io/ioutil"
	"log"
	"net/http"
	"os"
	// custom handlers
	"github.com/callistom/go-graphql-auth/authentication"
	"github.com/callistom/go-graphql-auth/handler"
	"github.com/callistom/go-graphql-auth/migrations"
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
		log.Fatal(err)
	}

	// parse schema file
	schema, err = graphql.ParseSchema(string(schemaFile), &handler.Resolver{})

	if err != nil {
		log.Fatal(err)
	}

}

func main() {

	// set logger
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, err := ioutil.ReadFile("graphiql.html")

		if err != nil {
			log.Fatal(err)
		}

		w.Write(page)
	}))

	// open db
	db, err = gorm.Open("postgres", "postgres://localhost/graphql?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	// start migrations
	migrated, migrationErr := migrations.CreateMigrations(db)

	if migrationErr != nil {
		logger.Printf("Migrations have been unsuccessful")
	}

	if migrated == true {
		logger.Printf("Migrations have been successful")
	}

	// close db last
	defer db.Close()

	// set terminal listner
	logger.Printf("Server is starting...")

	// set endpoint + authentication handler
	http.Handle("/graphql", authentication.Auth(&relay.Handler{Schema: schema}))

	// set terminal listner
	logger.Printf("Listening at http://localhost:8080")

	// listen and serve on correct url
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
