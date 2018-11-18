package main

import (
	"encoding/json"
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

func (mdb *mockDB) AddUser(id uint32, email string, name string, is_admin bool) error {
	if mdb.addedUsers == nil {
		mdb.addedUsers = make([]*models.User, 0)
	}
	mdb.addedUsers = append(mdb.addedUsers, &models.User{
		ID:      id,
		Email:   email,
		Name:    name,
		IsAdmin: is_admin,
	})
	return nil
}

func (mdb *mockDB) GetAllVisitedPaths() ([]*models.VisitedPath, error) {
	vps := make([]*models.VisitedPath, 0)
	vps = append(vps, &models.VisitedPath{
		Path: "/path1",
		Date: time.Date(2018, time.November, 17, 0, 0, 0, 0, time.UTC),
	})
	vps = append(vps, &models.VisitedPath{
		Path: "/path2",
		Date: time.Date(2018, time.November, 16, 0, 0, 0, 0, time.UTC),
	})
	return vps, nil
}

func (mdb *mockDB) AddVisitedPath(p string, ti time.Time) error {
	if mdb.addedVPs == nil {
		mdb.addedVPs = make([]*models.VisitedPath, 0)
	}
	mdb.addedVPs = append(mdb.addedVPs, &models.VisitedPath{
		Path: p,
		Date: ti,
	})
	return nil
}

// test handlers
func TestCanGetRootHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/abc", nil)

	db := &mockDB{}
	env := Env{db: db}
	http.HandlerFunc(env.rootHandler).ServeHTTP(rec, req)

	// check that the correct JSON strings were returned
	// read back in as string-string map
	var strs map[string]string
	err := json.Unmarshal([]byte(rec.Body.String()), &strs)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// check for expected length and values
	if len(strs) != 2 {
		t.Fatalf("expected len %d, got %d", 2, len(strs))
	}
	if strs["path"] != "/abc" {
		t.Errorf("expected %s, got %s", "/abc", strs["path"])
	}
	// don't check for exact date, b/c it'll vary per call
	// just make sure it exists and is non-nil
	if strs["date"] == "" {
		t.Errorf("expected non-nil date, got nil")
	}

	// and check that AddVisitedPath was called
	if len(db.addedVPs) != 1 {
		t.Errorf("expected len %d, got %d", 1, len(db.addedVPs))
	}
	if db.addedVPs[0].Path != "/abc" {
		t.Errorf("expected %s, got %s", "/abc", db.addedVPs[0].Path)
	}
}

func TestCannotPostRootHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/abc", nil)

	db := &mockDB{}
	env := Env{db: db}
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
	req, _ := http.NewRequest("GET", "/history", nil)

	db := &mockDB{}
	env := Env{db: db}
	http.HandlerFunc(env.historyHandler).ServeHTTP(rec, req)

	// check that the correct JSON strings were returned
	// read back in as string-string map
	var vals []map[string]string
	err := json.Unmarshal([]byte(rec.Body.String()), &vals)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// check for expected length and values
	if len(vals) != 2 {
		t.Fatalf("expected len %d, got %d", 2, len(vals))
	}

	path1 := vals[0]
	if len(path1) != 2 {
		t.Fatalf("expected len %d, got %d", 2, len(path1))
	}
	if path1["path"] != "/path1" {
		t.Errorf("expected %s, got %s", "/path1", path1["path"])
	}
	if path1["date"] != "2018-11-17T00:00:00Z" {
		t.Errorf("expected %s, got %s", "2018-11-17T00:00:00Z", path1["date"])
	}

	path2 := vals[1]
	if len(path2) != 2 {
		t.Fatalf("expected len %d, got %d", 2, len(path2))
	}
	if path2["path"] != "/path2" {
		t.Errorf("expected %s, got %s", "/path2", path2["path"])
	}
	if path2["date"] != "2018-11-16T00:00:00Z" {
		t.Errorf("expected %s, got %s", "2018-11-16T00:00:00Z", path2["date"])
	}
}

func TestCannotPostHistoryHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/history", nil)

	db := &mockDB{}
	env := Env{db: db}
	http.HandlerFunc(env.rootHandler).ServeHTTP(rec, req)

	// check that we got a 405
	if 405 != rec.Code {
		t.Errorf("Expected %d, got %d", 405, rec.Code)
	}
}

func TestIgnoreHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/favicon.ico", nil)

	db := &mockDB{}
	env := Env{db: db}
	http.HandlerFunc(env.ignoreHandler).ServeHTTP(rec, req)

	// check that we got a 404
	if 404 != rec.Code {
		t.Errorf("Expected %d, got %d", 404, rec.Code)
	}
}
