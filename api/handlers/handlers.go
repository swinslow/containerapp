package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/swinslow/containerapp/api/models"
)

// RegisterHandlers registers the api handler endpoints with the
// specified router, for the given environment.
func (env *Env) RegisterHandlers(router *mux.Router) {
	router.HandleFunc("/favicon.ico", env.ignoreHandler).Methods("GET")
	router.HandleFunc("/oauth/getToken", env.createTokenHandler).Methods("POST")
	router.HandleFunc("/landing", env.validateTokenMiddleware(env.landingHandler)).Methods("GET")
	router.HandleFunc("/admin/history", env.validateTokenMiddleware(env.historyHandler)).Methods("GET")
	router.HandleFunc("/admin/users", env.validateTokenMiddleware(env.getUsersHandler)).Methods("GET")
	router.HandleFunc("/admin/users", env.validateTokenMiddleware(env.newUserHandler)).Methods("POST")
	router.HandleFunc("/{rest:.*}", env.validateTokenMiddleware(env.rootHandler)).Methods("GET")
}

func (env *Env) rootHandler(w http.ResponseWriter, r *http.Request) {
	// send JSON responses
	w.Header().Set("Content-Type", "application/json")

	// we only take GET requests
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	d := time.Now()

	// pull User from context
	ctxCheck := r.Context().Value(userContextKey(0))
	if ctxCheck == nil {
		w.Header().Set("WWW-Authenticate", "Bearer")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "Authorization header with valid Bearer token required"}`)
		return
	}
	user := ctxCheck.(*models.User)
	if user.ID == 0 {
		w.Header().Set("WWW-Authenticate", "Bearer")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "unknown user %s"}`, user.Email)
		return
	}

	// authenticated, so proceed
	vp := models.VisitedPath{
		Path:   r.URL.Path,
		Date:   d,
		UserID: user.ID,
	}
	js, err := json.Marshal(&vp)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	fmt.Fprintf(w, string(js))

	// and add this one for future visits
	err = env.db.AddVisitedPath(r.URL.Path, d, user.ID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
}

func (env *Env) ignoreHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(404), 404)
}

func (env *Env) landingHandler(w http.ResponseWriter, r *http.Request) {
	// send JSON responses
	w.Header().Set("Content-Type", "application/json")

	// we only take GET requests
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	// pull User from context
	ctxCheck := r.Context().Value(userContextKey(0))
	if ctxCheck == nil {
		w.Header().Set("WWW-Authenticate", "Bearer")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "Authorization header with valid Bearer token required"}`)
		return
	}
	user := ctxCheck.(*models.User)

	// NOTE that we won't check here to see whether the token is for a
	// valid (i.e. known) user. As currently set up, anyone with an email
	// address can request a token and get to the /landing page. But
	// they won't be able to visit any other pages until they are
	// registered.
	js, err := json.Marshal(&user)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	fmt.Fprintf(w, string(js))
}
