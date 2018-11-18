package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/swinslow/containerapp/api/models"
)

// define mock Datastore

type mockDB struct {
	addedVPs []*models.VisitedPath
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

	// check that the correct text was returned
	expected := "Hello, path is /abc<br><br>\n"
	if expected != rec.Body.String() {
		t.Errorf("expected %v, got %v", expected, rec.Body.String())
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

	// check that the correct text was returned
	expected := `Previously visited paths:<br>
<ul>
<li>/path1 (2018-11-17 00:00:00 +0000 UTC)</li>
<li>/path2 (2018-11-16 00:00:00 +0000 UTC)</li>
</ul>
`
	if expected != rec.Body.String() {
		t.Errorf("expected %v, got %v", expected, rec.Body.String())
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
