package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/swinslow/containerapp/models"
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

	http.HandleFunc("/", env.rootHandler)
	fmt.Println("Listening on :" + WEBPORT)
	log.Fatal(http.ListenAndServe(":"+WEBPORT, nil))
}

// ===== environment setup =====

func setupEnv() (*Env, error) {
	db, err := models.NewDB("host=db sslmode=disable dbname=dev user=postgres-dev password=s3cr3tp4ssw0rd")
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

// ===== handlers =====

func (env *Env) rootHandler(w http.ResponseWriter, r *http.Request) {
	// we only take GET requests
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	fmt.Fprintf(w, "Hello, path is %s<br><br>\n", r.URL.Path)

	// get prior visited paths
	vpaths, err := env.db.GetAllVisitedPaths()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	fmt.Fprintf(w, "Previously visited paths:<br>\n<ul>\n")
	for _, vp := range vpaths {
		fmt.Fprintf(w, "<li>%s (%v)</li>\n", vp.Path, vp.Date)
	}
	fmt.Fprintf(w, "</ul>\n")

	// and add this one for future visits
	err = env.db.AddVisitedPath(r.URL.Path)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
}
