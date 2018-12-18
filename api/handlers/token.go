package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/swinslow/containerapp/api/models"
)

func (env *Env) githubLoginHandler(w http.ResponseWriter, r *http.Request) {
	// send JSON responses
	w.Header().Set("Content-Type", "application/json")

	// we only take GET requests
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	// get state query parameter
	states, ok := r.URL.Query()["state"]
	if !ok || len(states[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Must send a state query parameter"}`)
		return
	}
	state := states[0]

	// build and send OAuth redirect
	url := env.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// FIXME needs test suite added
func (env *Env) githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// first, retrieve the state string and keep it for the response
	state := r.FormValue("state")
	// FIXME decide whether it's an error if the state string is empty.
	// FIXME currently, assuming it's OK b/c it's up to the client whether
	// FIXME they care.

	// now, check the code
	code := r.FormValue("code")
	token, err := env.oauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Couldn't convert auth code to OAuth token"}`)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// use the code to get user data
	oauthClient := env.oauthConfig.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get(nil, "")
	if err != nil {
		fmt.Printf("client.Users.Get() failed with %v\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Printf("Logged in as GitHub user: %s\n", *user.Login)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (env *Env) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	// send JSON responses
	w.Header().Set("Content-Type", "application/json")

	// we only take POST requests
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	// get form values
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Couldn't parse form"}`)
		return
	}
	email := r.Form.Get("email")
	if email == "" {
		// don't use http.Error(), because we want to send custom
		// JSON text
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Must supply login email address in token request"}`)
		return
	}

	// create token using JWT secret key
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
	})
	tknString, err := tkn.SignedString([]byte(env.jwtSecretKey))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Couldn't create token"}`)
		return
	}

	// output as JSON
	fmt.Fprintf(w, `{"token": "%s"}`, tknString)
}

func sendAuthFail(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("WWW-Authenticate", "Bearer")
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprintf(w, `{"error": "Authorization header with valid Bearer token required"}`)
}

type userContextKey int

func (env *Env) validateTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// look for and extract the token from header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			sendAuthFail(w)
			return
		}

		// check that the auth header has the expected format
		// e.g. Authorization: Bearer ....
		if !strings.HasPrefix(authHeader, "Bearer ") {
			sendAuthFail(w)
			return
		}
		remainder := strings.TrimPrefix(authHeader, "Bearer ")

		// decrypt and validate the token
		token, err := jwt.Parse(remainder, func(tkn *jwt.Token) (interface{}, error) {
			// should this be HMAC when token was signed with HS256?
			if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Couldn't parse token method")
			}
			return []byte(env.jwtSecretKey), nil
		})
		if err != nil {
			sendAuthFail(w)
			return
		}

		var email string
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			email, ok = claims["email"].(string)
			if !ok {
				sendAuthFail(w)
				return
			}
		}

		// make sure this email also exists in the User database
		user, err := env.db.GetUserByEmail(email)
		if err != nil {
			user = &models.User{
				ID:      0,
				Email:   email,
				Name:    "",
				IsAdmin: false,
			}
		}

		// good to go! set context and move on
		ctx := r.Context()
		ctx = context.WithValue(ctx, userContextKey(0), user)
		next(w, r.WithContext(ctx))
	})
}
