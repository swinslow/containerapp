package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func (mdb *mockDB) GetUserByID(id uint32) (*models.User, error) {
	users, err := mdb.GetAllUsers()
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}
	// not found
	return nil, fmt.Errorf("user not found")
}

func (mdb *mockDB) GetUserByEmail(email string) (*models.User, error) {
	users, err := mdb.GetAllUsers()
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		if user.Email == email {
			return user, nil
		}
	}
	// not found
	return nil, fmt.Errorf("user not found")
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

// ===== helpers for tests

func confirmRecWasInvalidAuth(t *testing.T, rec *httptest.ResponseRecorder, es string) {
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

	// check that content type was application/json
	if header.Get("Content-Type") != "application/json" {
		t.Errorf("expected %v, got %v", "application/json", header.Get("Content-Type"))
	}

	// check that the right "error" JSON string was returned
	wantString := `{"error": "` + es + `"}`
	if rec.Body.String() != wantString {
		t.Fatalf("expected %s, got %s", wantString, rec.Body.String())
	}
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

	// add User to context (assumes validation has already occurred)
	user, err := db.GetUserByEmail("janedoe@example.com")
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}
	ctx := req.Context()
	ctx = context.WithValue(ctx, userContextKey(0), user)
	req = req.WithContext(ctx)
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
	// check for user ID from token
	if vpGot.UserID != user.ID {
		t.Errorf("expected %v, got %v", user.ID, vpGot.UserID)
	}

	// and check that AddVisitedPath was called
	if len(db.addedVPs) != 1 {
		t.Errorf("expected len %d, got %d", 1, len(db.addedVPs))
	}
	if db.addedVPs[0].Path != "/abc" {
		t.Errorf("expected %v, got %v", "/abc", db.addedVPs[0].Path)
	}
}

func TestCannotGetRootHandlerWithoutValidUserInContext(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/abc", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	db := &mockDB{}
	env := Env{db: db, jwtSecretKey: "keyForTesting"}

	// add User with ID 0 to context (unknown user)
	user := &models.User{
		ID:      0,
		Email:   "unknown@example.com",
		Name:    "",
		IsAdmin: false,
	}
	ctx := req.Context()
	ctx = context.WithValue(ctx, userContextKey(0), user)
	req = req.WithContext(ctx)
	http.HandlerFunc(env.rootHandler).ServeHTTP(rec, req)

	confirmRecWasInvalidAuth(t, rec, "unknown user unknown@example.com")

	// and check that AddVisitedPath was not called
	if len(db.addedVPs) != 0 {
		t.Errorf("expected len %d, got %d", 0, len(db.addedVPs))
	}
}

func TestCannotGetRootHandlerWithNoUserInContext(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/abc", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	db := &mockDB{}
	env := Env{db: db, jwtSecretKey: "keyForTesting"}

	// not adding any User to context
	http.HandlerFunc(env.rootHandler).ServeHTTP(rec, req)

	confirmRecWasInvalidAuth(t, rec, "Authorization header with valid Bearer token required")

	// and check that AddVisitedPath was not called
	if len(db.addedVPs) != 0 {
		t.Errorf("expected len %d, got %d", 0, len(db.addedVPs))
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

// ===== Route: GET /landing =====

func TestCanGetLandingHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/landing", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	db := &mockDB{}
	env := Env{db: db, jwtSecretKey: "keyForTesting"}

	// add User to context (assumes validation has already occurred)
	user, err := db.GetUserByEmail("janedoe@example.com")
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}
	ctx := req.Context()
	ctx = context.WithValue(ctx, userContextKey(0), user)
	req = req.WithContext(ctx)
	http.HandlerFunc(env.landingHandler).ServeHTTP(rec, req)

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
	gotUser := &models.User{}
	err = json.Unmarshal([]byte(rec.Body.String()), &gotUser)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	if gotUser.ID != 914611345 {
		t.Errorf("expected %v, got %v", 914611345, gotUser.ID)
	}
	if gotUser.Email != "janedoe@example.com" {
		t.Errorf("expected %v, got %v", "janedoe@example.com", gotUser.Email)
	}
	if gotUser.Name != "Jane Doe" {
		t.Errorf("expected %v, got %v", "Jane Doe", gotUser.Name)
	}
	if gotUser.IsAdmin != true {
		t.Errorf("expected %v, got %v", true, gotUser.IsAdmin)
	}
}

func TestCannotPostLandingHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/landing", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	db := &mockDB{}
	env := Env{db: db, jwtSecretKey: "keyForTesting"}
	http.HandlerFunc(env.landingHandler).ServeHTTP(rec, req)

	// check that we got a 405
	if 405 != rec.Code {
		t.Errorf("Expected %d, got %d", 405, rec.Code)
	}
}
