package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/swinslow/containerapp/api/handlers"
)

func main() {
	var WEBPORT string
	if WEBPORT = os.Getenv("WEBPORT"); WEBPORT == "" {
		WEBPORT = "3001"
	}

	// set up database object and environment
	env, err := handlers.SetupEnv()
	if err != nil {
		log.Panic(err)
	}

	// create router and register handlers
	router := mux.NewRouter()

	env.RegisterHandlers(router)

	// set up CORS
	headers := []string{"X-Requested-With", "Content-Type", "Authorization"}
	methods := []string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}
	origins := []string{"http://localhost:3000"}
	cors := gh.CORS(
		gh.AllowedHeaders(headers),
		gh.AllowedMethods(methods),
		gh.AllowedOrigins(origins))

	fmt.Println("Listening on :" + WEBPORT)
	log.Fatal(http.ListenAndServe(":"+WEBPORT, cors(router)))
}
