package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/swinslow/containerapp/api/models"
)

// Env is the environment for the web handlers.
type Env struct {
	db models.Datastore
}

func main() {
	var WEBPORT string
	if WEBPORT = os.Getenv("WEBPORT"); WEBPORT == "" {
		WEBPORT = "3001"
	}

	// set up database object and environment
	env, err := setupEnv()
	if err != nil {
		log.Panic(err)
	}

	// create router and register handlers
	router := mux.NewRouter()

	registerHandlers(router, env)

	// set up CORS
	headers := []string{"X-Requested-With", "Content-Type", "Authorization"}
	methods := []string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}
	origins := []string{"http://localhost:3000"}

	fmt.Println("Listening on :" + WEBPORT)
	log.Fatal(http.ListenAndServe(":"+WEBPORT, handlers.CORS(
		handlers.AllowedHeaders(headers),
		handlers.AllowedMethods(methods),
		handlers.AllowedOrigins(origins))(router)))
}

// ===== environment setup =====

func setupEnv() (*Env, error) {
	db, err := models.NewDB("host=db sslmode=disable dbname=dev user=postgres-dev")
	if err != nil {
		return nil, err
	}

	err = db.InitDBTables()
	if err != nil {
		return nil, err
	}

	env := &Env{db}
	return env, nil
}
