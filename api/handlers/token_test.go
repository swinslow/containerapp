package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/swinslow/containerapp/api/models"
)

func TestCanPostCreateTokenHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	// send as URL-encoded form data b/c that's what'll happen
	// in the OAuth token workflow
	// for now, though, we'll just trust whatever email address they send
	data := url.Values{}
	data.Set("email", "janedoe@example.com")
	req, err := http.NewRequest("POST", "/oauth/getToken", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}
	// also have to set this header so that Form can get populated by the server
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	db := &mockDB{}
	env := Env{db: db, jwtSecretKey: "keyForTesting"}
	http.HandlerFunc(env.createTokenHandler).ServeHTTP(rec, req)

	// check that we got a 200 (OK)
	if 200 != rec.Code {
		t.Errorf("Expected %d, got %d", 200, rec.Code)
	}

	// check that content type was application/json
	header := rec.Result().Header
	if header.Get("Content-Type") != "application/json" {
		t.Errorf("expected %v, got %v", "application/json", header.Get("Content-Type"))
	}

	// check that a JSON response with a JWT was returned
	// (we won't currently try to decrypt it, just confirm it exists)
	rj := map[string]string{}
	err = json.Unmarshal([]byte(rec.Body.String()), &rj)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	if len(rj) != 1 {
		t.Fatalf("expected len %d, got %d", 1, len(rj))
	}

	if _, ok := rj["token"]; !ok {
		t.Errorf("expected token, got no token key")
	}
}

func TestCannotGetCreateTokenHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/oauth/getToken", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	db := &mockDB{}
	env := Env{db: db, jwtSecretKey: "keyForTesting"}
	http.HandlerFunc(env.createTokenHandler).ServeHTTP(rec, req)

	// check that we got a 405
	if 405 != rec.Code {
		t.Errorf("Expected %d, got %d", 405, rec.Code)
	}
}

func TestCannotPostCreateTokenHandlerWithEmptyEmailString(t *testing.T) {
	rec := httptest.NewRecorder()
	data := url.Values{}
	data.Set("email", "")
	req, err := http.NewRequest("POST", "/oauth/getToken", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	db := &mockDB{}
	env := Env{db: db, jwtSecretKey: "keyForTesting"}
	http.HandlerFunc(env.createTokenHandler).ServeHTTP(rec, req)

	// check that we got a 400 (Bad Request)
	if 400 != rec.Code {
		t.Errorf("Expected %d, got %d", 400, rec.Code)
	}

	// check that content type was application/json
	header := rec.Result().Header
	if header.Get("Content-Type") != "application/json" {
		t.Errorf("expected %v, got %v", "application/json", header.Get("Content-Type"))
	}

	// check that a JSON response with an error was returned
	rj := map[string]string{}
	err = json.Unmarshal([]byte(rec.Body.String()), &rj)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	if len(rj) != 1 {
		t.Fatalf("expected len %d, got %d", 1, len(rj))
	}

	if _, ok := rj["error"]; !ok {
		t.Errorf("expected error, got no error key")
	}
}

func TestCannotPostCreateTokenHandlerWithoutEmailValue(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/oauth/getToken", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	db := &mockDB{}
	env := Env{db: db, jwtSecretKey: "keyForTesting"}
	http.HandlerFunc(env.createTokenHandler).ServeHTTP(rec, req)

	// check that we got a 400 (Bad Request)
	if 400 != rec.Code {
		t.Errorf("Expected %d, got %d", 400, rec.Code)
	}

	// check that content type was application/json
	header := rec.Result().Header
	if header.Get("Content-Type") != "application/json" {
		t.Errorf("expected %v, got %v", "application/json", header.Get("Content-Type"))
	}

	// check that a JSON response with an error was returned
	rj := map[string]string{}
	err = json.Unmarshal([]byte(rec.Body.String()), &rj)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	if len(rj) != 1 {
		t.Fatalf("expected len %d, got %d", 1, len(rj))
	}

	if _, ok := rj["error"]; !ok {
		t.Errorf("expected error, got no error key")
	}
}

// ===== test middleware =====

// sample handler for testing middleware
func (env *Env) testHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if user := r.Context().Value(userContextKey(0)).(*models.User); user != nil {
		fmt.Fprintf(w, `{"email": "%s", "id": %d}`, user.Email, user.ID)
	} else {
		fmt.Fprintf(w, `{"error": "couldn't get context"}`)
	}
}

func TestCanValidateTokenMiddleware(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/testRoute", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// create token with testing key and set header
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": "janedoe@example.com",
	})
	tokenString, err := token.SignedString([]byte("keyForTesting"))
	if err != nil {
		t.Fatalf("couldn't create token for testing: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+tokenString)

	db := &mockDB{}
	env := Env{db: db, jwtSecretKey: "keyForTesting"}
	wrappedHandler := env.validateTokenMiddleware(env.testHandler)
	http.HandlerFunc(wrappedHandler).ServeHTTP(rec, req)

	// check that we got a 200 (OK)
	if 200 != rec.Code {
		t.Errorf("Expected %d, got %d", 200, rec.Code)
	}

	// check that content type was application/json
	header := rec.Result().Header
	if header.Get("Content-Type") != "application/json" {
		t.Errorf("expected %v, got %v", "application/json", header.Get("Content-Type"))
	}

	wantString := `{"email": "janedoe@example.com", "id": 914611345}`
	if rec.Body.String() != wantString {
		t.Errorf("expected %s, got %s", wantString, rec.Body.String())
	}
}

func TestCanValidateTokenMiddlewareForUnknownEmailButIDIsZero(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/testRoute", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// create token with testing key and set header
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": "unknown@example.com",
	})
	tokenString, err := token.SignedString([]byte("keyForTesting"))
	if err != nil {
		t.Fatalf("couldn't create token for testing: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+tokenString)

	db := &mockDB{}
	env := Env{db: db, jwtSecretKey: "keyForTesting"}
	wrappedHandler := env.validateTokenMiddleware(env.testHandler)
	http.HandlerFunc(wrappedHandler).ServeHTTP(rec, req)

	// check that we got a 200 (OK)
	if 200 != rec.Code {
		t.Errorf("Expected %d, got %d", 200, rec.Code)
	}

	// check that content type was application/json
	header := rec.Result().Header
	if header.Get("Content-Type") != "application/json" {
		t.Errorf("expected %v, got %v", "application/json", header.Get("Content-Type"))
	}

	wantString := `{"email": "unknown@example.com", "id": 0}`
	if rec.Body.String() != wantString {
		t.Errorf("expected %s, got %s", wantString, rec.Body.String())
	}
}

func TestCannotValidateTokenWithNoAuthHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/testRoute", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	db := &mockDB{}
	env := Env{db: db, jwtSecretKey: "keyForTesting"}
	wrappedHandler := env.validateTokenMiddleware(env.testHandler)
	http.HandlerFunc(wrappedHandler).ServeHTTP(rec, req)

	// check that we got a 401 (Unauthorized)
	if 401 != rec.Code {
		t.Errorf("Expected %d, got %d", 401, rec.Code)
	}

	// check that we got a WWW-Authenticate header
	// (see https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/401)
	header := rec.Result().Header
	wantHeader := "Bearer"
	gotHeader := header.Get("WWW-Authenticate")
	if gotHeader != wantHeader {
		t.Errorf("expected %v, got %v", wantHeader, gotHeader)
	}

	wantBody := `{"error": "Authentication header with valid Bearer token required"}`
	if rec.Body.String() != wantBody {
		t.Errorf("expected %s, got %s", wantBody, rec.Body.String())
	}
}

func TestCannotValidateTokenWithInvalidAuthHeaders(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/testRoute", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// set an invalid JWT token value
	req.Header.Set("Authorization", "Bearer BLAHinvalid")

	db := &mockDB{}
	env := Env{db: db, jwtSecretKey: "keyForTesting"}
	wrappedHandler := env.validateTokenMiddleware(env.testHandler)
	http.HandlerFunc(wrappedHandler).ServeHTTP(rec, req)

	// check that we got a 401 (Unauthorized)
	if 401 != rec.Code {
		t.Errorf("Expected %d, got %d", 401, rec.Code)
	}

	// check that we got a WWW-Authenticate header
	// (see https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/401)
	header := rec.Result().Header
	wantHeader := "Bearer"
	gotHeader := header.Get("WWW-Authenticate")
	if gotHeader != wantHeader {
		t.Errorf("expected %v, got %v", wantHeader, gotHeader)
	}

	wantBody := `{"error": "Authentication header with valid Bearer token required"}`
	if rec.Body.String() != wantBody {
		t.Errorf("expected %s, got %s", wantBody, rec.Body.String())
	}
}

func TestCannotValidateTokenWithNoBearerInHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/testRoute", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// create token with testing key and set header
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": "janedoe@example.com",
	})
	tokenString, err := token.SignedString([]byte("keyForTesting"))
	if err != nil {
		t.Fatalf("couldn't create token for testing: %v", err)
	}
	req.Header.Set("Authorization", tokenString)

	db := &mockDB{}
	env := Env{db: db, jwtSecretKey: "keyForTesting"}
	wrappedHandler := env.validateTokenMiddleware(env.testHandler)
	http.HandlerFunc(wrappedHandler).ServeHTTP(rec, req)

	// check that we got a 401 (Unauthorized)
	if 401 != rec.Code {
		t.Errorf("Expected %d, got %d", 401, rec.Code)
	}

	// check that we got a WWW-Authenticate header
	// (see https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/401)
	header := rec.Result().Header
	wantHeader := "Bearer"
	gotHeader := header.Get("WWW-Authenticate")
	if gotHeader != wantHeader {
		t.Errorf("expected %v, got %v", wantHeader, gotHeader)
	}

	wantBody := `{"error": "Authentication header with valid Bearer token required"}`
	if rec.Body.String() != wantBody {
		t.Errorf("expected %s, got %s", wantBody, rec.Body.String())
	}
}
