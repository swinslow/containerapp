package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/swinslow/containerapp/api/models"
)

// extractAdminUser takes a Request and extracts and returns the User object,
// confirming that it is an admin User. If the User cannot be extracted,
// is unknown, or is not an admin, the appropriate HTTP headers and content
// are written and nil is returned (presumably to flag that no further
// processing should occur).
func extractAdminUser(w http.ResponseWriter, r *http.Request) *models.User {
	// pull User from context
	ctxCheck := r.Context().Value(userContextKey(0))
	if ctxCheck == nil {
		w.Header().Set("WWW-Authenticate", "Bearer")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "Authorization header with valid Bearer token required"}`)
		return nil
	}
	user := ctxCheck.(*models.User)
	if user.ID == 0 {
		w.Header().Set("WWW-Authenticate", "Bearer")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "unknown user %s"}`, user.Email)
		return nil
	}

	// admin access required for this resource
	if !user.IsAdmin {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, `{"error": "admin access required"}`)
		return nil
	}

	return user
}

func (env *Env) historyHandler(w http.ResponseWriter, r *http.Request) {
	// send JSON responses
	w.Header().Set("Content-Type", "application/json")

	// we only take GET requests
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	user := extractAdminUser(w, r)
	if user == nil {
		return
	}

	// get prior visited paths
	vpaths, err := env.db.GetAllVisitedPaths()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// output as JSON
	js, err := json.Marshal(vpaths)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	fmt.Fprintf(w, string(js))
}

func (env *Env) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	// send JSON responses
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	user := extractAdminUser(w, r)
	if user == nil {
		return
	}

	// once we're here, we're talking to an admin user

	// return list of all users
	users, err := env.db.GetAllUsers()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// output as JSON
	js, err := json.Marshal(users)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	fmt.Fprintf(w, string(js))
}

type newUserReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (env *Env) newUserHandler(w http.ResponseWriter, r *http.Request) {
	setToAdmin := false

	// send JSON responses
	w.Header().Set("Content-Type", "application/json")

	// we only take POST requests
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	// see if this is a brand new, uninitialized database -- in other
	// words, if it has no users. if so, we'll create the first user
	// even though there's no auth (because they can't auth yet!)
	// FIXME in a production system, this is a dumb way to pick the
	// FIXME first admin user.
	if users, err := env.db.GetAllUsers(); err == nil && len(users) == 0 {
		// new database, so skip the admin check.
		// also flag that the new user should be an admin.
		setToAdmin = true
	} else {
		user := extractAdminUser(w, r)
		if user == nil {
			return
		}
	}

	// once we're here, we're talking to an admin user

	// extract JSON content
	var newUser newUserReq
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// make sure this user doesn't already exist in database
	if userCheck, err := env.db.GetUserByEmail(newUser.Email); err == nil && userCheck != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "user with email %s already exists"}`, newUser.Email)
		return
	}

	// generate random new, non-zero, currently-unused ID within permitted range
	// FIXME note that this process for finding an unused ID will get decreasingly
	// FIXME usable if you have a massive number of users.
	var newID uint32
	for newID == 0 {
		testID := rand.Int31()
		// make it a positive, non-negative number
		if testID == 0 {
			continue
		}
		if testID < 0 {
			testID = -testID
		}
		utestID := uint32(testID)
		// check this ID doesn't already exist in database
		if userCheck, err := env.db.GetUserByID(utestID); err == nil && userCheck != nil {
			continue
		}
		// if we get here, it's good to use
		newID = utestID
	}

	// create and save new user
	// FIXME note that the new ID selection is NOT safe across multiple
	// FIXME requests. two separate requests could both claim the same
	// FIXME ID as available above, and then both try to save them here.
	// FIXME This will be prevented by the database, presumably, but it
	// FIXME should be addressed in a real production system.
	err = env.db.AddUser(newID, newUser.Email, newUser.Name, setToAdmin)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "server error saving new user, please check values and try again"}`)
		return
	}

	// success!
	finalUser := models.User{
		ID:      newID,
		Email:   newUser.Email,
		Name:    newUser.Name,
		IsAdmin: setToAdmin,
	}
	w.WriteHeader(http.StatusCreated)
	js, err := json.Marshal(finalUser)
	if err != nil {
		// created resource, but some problem with marshalling the JSON
		// response; just return empty in this case
		return
	}
	fmt.Fprintf(w, string(js))
}
