package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func registerHandlers(router *mux.Router, env *Env) {
	router.HandleFunc("/history", env.historyHandler).Methods("GET")
	router.HandleFunc("/{rest:.*}", env.rootHandler).Methods("GET")
}

func (env *Env) rootHandler(w http.ResponseWriter, r *http.Request) {
	// we only take GET requests
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	fmt.Fprintf(w, "Hello, path is %s<br><br>\n", r.URL.Path)

	// and add this one for future visits
	err := env.db.AddVisitedPath(r.URL.Path, time.Now())
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
}

func (env *Env) historyHandler(w http.ResponseWriter, r *http.Request) {
	// we only take GET requests
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

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
}
