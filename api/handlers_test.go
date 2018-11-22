package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/swinslow/containerapp/api/models"
)

// define mock Datastore

type mockDB struct {
	addedVPs   []*models.VisitedPath
	addedUsers []*models.User
}

func (mdb *mockDB) GetAllUsers() ([]*models.User, error) {
	users := make([]*models.User, 0)
	users = append(users, &models.User{
		ID:      91461,
		Email:   "johndoe@example.com",
		Name:    "John Doe",
		IsAdmin: false,
	})
	users = append(users, &models.User{
		ID:      914611345,
		Email:   "janedoe@example.com",
		Name:    "Jane Doe",
		IsAdmin: true,
	})
	return users, nil
}

func (mdb *mockDB) AddUser(id uint32, email string, name string, isAdmin bool) error {
	if mdb.addedUsers == nil {
		mdb.addedUsers = make([]*models.User, 0)
	}
	mdb.addedUsers = append(mdb.addedUsers, &models.User{
		ID:      id,
		Email:   email,
		Name:    name,
		IsAdmin: isAdmin,
	})
	return nil
}

func (mdb *mockDB) GetAllVisitedPaths() ([]*models.VisitedPath, error) {
	vps := make([]*models.VisitedPath, 0)
	vps = append(vps, &models.VisitedPath{
		Path:   "/path1",
		Date:   time.Date(2018, time.November, 17, 0, 0, 0, 0, time.UTC),
		UserID: 49185,
	})
	vps = append(vps, &models.VisitedPath{
		Path:   "/path2",
		Date:   time.Date(2018, time.November, 16, 0, 0, 0, 0, time.UTC),
		UserID: 847102,
	})
	return vps, nil
}

func (mdb *mockDB) GetAllVisitedPathsForUserID(uint32) ([]*models.VisitedPath, error) {
	vps := make([]*models.VisitedPath, 0)
	vps = append(vps, &models.VisitedPath{
		Path:   "/path1",
		Date:   time.Date(2018, time.November, 17, 0, 0, 0, 0, time.UTC),
		UserID: 49185,
	})
	return vps, nil
}

func (mdb *mockDB) AddVisitedPath(p string, ti time.Time, userID uint32) error {
	if mdb.addedVPs == nil {
		mdb.addedVPs = make([]*models.VisitedPath, 0)
	}
	mdb.addedVPs = append(mdb.addedVPs, &models.VisitedPath{
		Path:   p,
		Date:   ti,
		UserID: userID,
	})
	return nil
}

// ===== test handlers =====

func TestCanGetRootHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/abc", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	db := &mockDB{}
	env := Env{db: db, jwtSecretKey: "keyForTesting"}
	http.HandlerFunc(env.rootHandler).ServeHTTP(rec, req)

	// check that we got a 200 (OK)
	if 200 != rec.Code {
		t.Errorf("Expected %d, got %d", 200, rec.Code)
	}

	// check that content type was application/json
	header := rec.Result().Header
	if header.Get("Content-Type") != "application/json" {
		t.Errorf("expected %v, got %v", "application/json", header.Get("Content-Type"))
	}

	// check that the correct JSON strings were returned
	vpGot := &models.VisitedPath{}
	err = json.Unmarshal([]byte(rec.Body.String()), &vpGot)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// check for expected values
	if vpGot.Path != "/abc" {
		t.Errorf("expected %v, got %v", "/abc", vpGot.Path)
	}
	// don't check for exact date, b/c it'll vary per call
	// FIXME for now, checking for hard-coded user id
	if vpGot.UserID != 1001 {
		t.Errorf("expected %v, got %v", 1001, vpGot.UserID)
	}

	// and check that AddVisitedPath was called
	if len(db.addedVPs) != 1 {
		t.Errorf("expected len %d, got %d", 1, len(db.addedVPs))
	}
	if db.addedVPs[0].Path != "/abc" {
		t.Errorf("expected %v, got %v", "/abc", db.addedVPs[0].Path)
	}
}

func TestCannotPostRootHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/abc", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	db := &mockDB{}
	env := Env{db: db, jwtSecretKey: "keyForTesting"}
	http.HandlerFunc(env.rootHandler).ServeHTTP(rec, req)

	// check that we got a 405
	if 405 != rec.Code {
		t.Errorf("Expected %d, got %d", 405, rec.Code)
	}

	// and check that AddVisitedPath was not called
	if len(db.addedVPs) != 0 {
		t.Errorf("expected len %d, got %d", 0, len(db.addedVPs))
	}
}

func TestCanGetHistory(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/history", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	db := &mockDB{}
	env := Env{db: db, jwtSecretKey: "keyForTesting"}
	http.HandlerFunc(env.historyHandler).ServeHTTP(rec, req)

	// check that we got a 200 (OK)
	if 200 != rec.Code {
		t.Errorf("Expected %d, got %d", 200, rec.Code)
	}

	// check that content type was application/json
	header := rec.Result().Header
	if header.Get("Content-Type") != "application/json" {
		t.Errorf("expected %v, got %v", "application/json", header.Get("Content-Type"))
	}

	// check that the correct JSON strings were returned
	// read back in as slice of VisitedPaths
	var vals []*models.VisitedPath
	err = json.Unmarshal([]byte(rec.Body.String()), &vals)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// check for expected length and values
	if len(vals) != 2 {
		t.Fatalf("expected len %d, got %d", 2, len(vals))
	}

	vp1 := vals[0]
	if vp1.Path != "/path1" {
		t.Errorf("expected %v, got %v", "/vp1", vp1.Path)
	}
	wantDate1 := time.Date(2018, time.November, 17, 0, 0, 0, 0, time.UTC)
	if vp1.Date != wantDate1 {
		t.Errorf("expected %v, got %v", wantDate1, vp1.Date)
	}
	if vp1.UserID != 49185 {
		t.Errorf("expected %v, got %v", 49185, vp1.UserID)
	}

	vp2 := vals[1]
	if vp2.Path != "/path2" {
		t.Errorf("expected %v, got %v", "/vp2", vp2.Path)
	}
	wantDate2 := time.Date(2018, time.November, 16, 0, 0, 0, 0, time.UTC)
	if vp2.Date != wantDate2 {
		t.Errorf("expected %v, got %v", wantDate2, vp2.Date)
	}
	if vp2.UserID != 847102 {
		t.Errorf("expected %v, got %v", 847102, vp2.UserID)
	}
}

func TestCannotPostHistoryHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/history", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	db := &mockDB{}
	env := Env{db: db, jwtSecretKey: "keyForTesting"}
	http.HandlerFunc(env.rootHandler).ServeHTTP(rec, req)

	// check that we got a 405
	if 405 != rec.Code {
		t.Errorf("Expected %d, got %d", 405, rec.Code)
	}
}

func TestIgnoreHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/favicon.ico", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	db := &mockDB{}
	env := Env{db: db, jwtSecretKey: "keyForTesting"}
	http.HandlerFunc(env.ignoreHandler).ServeHTTP(rec, req)

	// check that we got a 404
	if 404 != rec.Code {
		t.Errorf("Expected %d, got %d", 404, rec.Code)
	}
}

func TestCanPostCreateTokenHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	// send as URL-encoded form data b/c that's what'll happen
	// in the OAuth token workflow
	// for now, though, we'll just trust whatever email address they send
	data := url.Values{}
	data.Set("email", "janedoe@example.com")
	req, err := http.NewRequest("POST", "/getToken", strings.NewReader(data.Encode()))
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
	req, err := http.NewRequest("GET", "/getToken", nil)
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
	req, err := http.NewRequest("POST", "/getToken", strings.NewReader(data.Encode()))
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
	req, err := http.NewRequest("POST", "/getToken", nil)
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
	if userEmail := r.Context().Value(emailContextKey(0)).(string); userEmail != "" {
		fmt.Fprintf(w, "got %s from context", userEmail)
	} else {
		fmt.Fprintf(w, `couldn't get context`)
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

	wantString := "got janedoe@example.com from context"
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
