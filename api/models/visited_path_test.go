package models

import (
	"encoding/json"
	"testing"
	"time"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestShouldGetAllVisitedPaths(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	goodbyeDate := time.Date(2018, time.November, 17, 0, 0, 0, 0, time.UTC)
	goneDate := time.Date(2018, time.November, 16, 0, 0, 0, 0, time.UTC)
	helloDate := time.Date(2018, time.November, 15, 0, 0, 0, 0, time.UTC)

	sentRows := sqlmock.NewRows([]string{"path", "visit_date"}).
		AddRow("/goodbye", goodbyeDate).
		AddRow("/gone", goneDate).
		AddRow("/hello", helloDate)
	mock.ExpectQuery("SELECT path, visit_date FROM visitedpaths ORDER BY visit_date DESC").WillReturnRows(sentRows)

	// run the tested function
	gotRows, err := db.GetAllVisitedPaths()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values
	if len(gotRows) != 3 {
		t.Fatalf("expected len %d, got %d", 3, len(gotRows))
	}
	vp0 := gotRows[0]
	if vp0.Path != "/goodbye" {
		t.Errorf("expected %s, got %s", "/goodbye", vp0.Path)
	}
	if vp0.Date != goodbyeDate {
		t.Errorf("expected %s, got %s", goodbyeDate, vp0.Date)
	}
	vp1 := gotRows[1]
	if vp1.Path != "/gone" {
		t.Errorf("expected %s, got %s", "/gone", vp1.Path)
	}
	if vp1.Date != goneDate {
		t.Errorf("expected %s, got %s", goneDate, vp1.Date)
	}
	vp2 := gotRows[2]
	if vp2.Path != "/hello" {
		t.Errorf("expected %s, got %s", "/hello", vp2.Path)
	}
	if vp2.Date != helloDate {
		t.Errorf("expected %s, got %s", helloDate, vp2.Date)
	}

}

func TestShouldAddVisitedPath(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	helloDate := time.Date(2018, time.November, 15, 0, 0, 0, 0, time.UTC)

	regexStmt := `[INSERT INTO visitedpaths(path, visit_date) VALUES (\$1, \$2)]`
	mock.ExpectPrepare(regexStmt)
	stmt := "INSERT INTO visitedpaths"
	mock.ExpectExec(stmt).
		WithArgs("hello", helloDate).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// run the tested function
	err = db.AddVisitedPath("hello", helloDate)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// JSON marshalling and unmarshalling
func TestVisitedPathCanMarshalToJSON(t *testing.T) {
	vp := &VisitedPath{
		Path: "/abc",
		Date: time.Date(2018, time.November, 17, 0, 0, 0, 0, time.UTC),
	}

	js, err := json.Marshal(vp)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// read back in as string-string map
	var strs map[string]string
	err = json.Unmarshal(js, &strs)
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
	if strs["date"] != "2018-11-17T00:00:00Z" {
		t.Errorf("expected %s, got %s", "2018-11-17T00:00:00Z", strs["date"])
	}
}

func TestVisitedPathCanUnmarshalFromJSON(t *testing.T) {
	vp := &VisitedPath{}
	js := []byte(`{"path":"/def", "date":"2018-11-17T20:43:00Z"}`)

	err := json.Unmarshal(js, vp)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// check values
	if vp.Path != "/def" {
		t.Errorf("expected %v, got %v", "/def", vp.Path)
	}
	wantDate := time.Date(2018, time.November, 17, 20, 43, 0, 0, time.UTC)
	if vp.Date != wantDate {
		t.Errorf("expected %v, got %v", wantDate, vp.Date)
	}
}
