package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/swinslow/containerapp/api/models"
)

func TestCanGetHistory(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/admin/history", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	db := &mockDB{}
	env := Env{db: db, jwtSecretKey: "keyForTesting"}
	// add admin User to context (assumes validation has already occurred)
	user, err := db.GetUserByEmail("janedoe@example.com")
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}
	ctx := req.Context()
	ctx = context.WithValue(ctx, userContextKey(0), user)
	req = req.WithContext(ctx)
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

func TestCannotGetHistoryWithNoUserInContext(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/admin/history", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	db := &mockDB{}
	env := Env{db: db, jwtSecretKey: "keyForTesting"}
	// not adding any User to context
	http.HandlerFunc(env.historyHandler).ServeHTTP(rec, req)

	confirmRecWasInvalidAuth(t, rec, "Authentication header with valid Bearer token required")
}

func TestCannotGetHistoryWithoutValidUserInContext(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/admin/history", nil)
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
	http.HandlerFunc(env.historyHandler).ServeHTTP(rec, req)

	confirmRecWasInvalidAuth(t, rec, "unknown user unknown@example.com")
}

func TestCannotGetHistoryWithoutAdmindUserInContext(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/admin/history", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	db := &mockDB{}
	env := Env{db: db, jwtSecretKey: "keyForTesting"}

	// add User with ID 0 to context (unknown user)
	user, err := db.GetUserByEmail("johndoe@example.com")
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}
	ctx := req.Context()
	ctx = context.WithValue(ctx, userContextKey(0), user)
	req = req.WithContext(ctx)
	http.HandlerFunc(env.historyHandler).ServeHTTP(rec, req)

	// check that we got a 403 (Forbidden)
	if 403 != rec.Code {
		t.Errorf("Expected %d, got %d", 403, rec.Code)
	}

	// check that content type was application/json
	header := rec.Result().Header
	if header.Get("Content-Type") != "application/json" {
		t.Errorf("expected %v, got %v", "application/json", header.Get("Content-Type"))
	}

	// check that the right "error" JSON string was returned
	wantString := `{"error": "admin access required"}`
	if rec.Body.String() != wantString {
		t.Fatalf("expected %s, got %s", wantString, rec.Body.String())
	}
}

func TestCannotPostHistoryHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/admin/history", nil)
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
