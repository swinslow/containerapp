package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/swinslow/containerapp/api/models"
)

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
	fmt.Fprintf(w, `{"error": "Authentication header with valid Bearer token required"}`)
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
