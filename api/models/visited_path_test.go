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
	goodbyeUserID := uint32(2918592)
	goneDate := time.Date(2018, time.November, 16, 0, 0, 0, 0, time.UTC)
	goneUserID := uint32(56)
	helloDate := time.Date(2018, time.November, 15, 0, 0, 0, 0, time.UTC)
	helloUserID := uint32(2918592)

	sentRows := sqlmock.NewRows([]string{"path", "visit_date", "user_id"}).
		AddRow("/goodbye", goodbyeDate, goodbyeUserID).
		AddRow("/gone", goneDate, goneUserID).
		AddRow("/hello", helloDate, helloUserID)
	mock.ExpectQuery("SELECT path, visit_date, user_id FROM visitedpaths ORDER BY visit_date DESC").WillReturnRows(sentRows)

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
		t.Errorf("expected %v, got %v", "/goodbye", vp0.Path)
	}
	if vp0.Date != goodbyeDate {
		t.Errorf("expected %v, got %v", goodbyeDate, vp0.Date)
	}
	if vp0.UserID != goodbyeUserID {
		t.Errorf("expected %v, got %v", goodbyeUserID, vp0.UserID)
	}
	vp1 := gotRows[1]
	if vp1.Path != "/gone" {
		t.Errorf("expected %v, got %v", "/gone", vp1.Path)
	}
	if vp1.Date != goneDate {
		t.Errorf("expected %v, got %v", goneDate, vp1.Date)
	}
	if vp1.UserID != goneUserID {
		t.Errorf("expected %v, got %v", goneUserID, vp1.UserID)
	}
	vp2 := gotRows[2]
	if vp2.Path != "/hello" {
		t.Errorf("expected %v, got %v", "/hello", vp2.Path)
	}
	if vp2.Date != helloDate {
		t.Errorf("expected %v, got %v", helloDate, vp2.Date)
	}
	if vp2.UserID != helloUserID {
		t.Errorf("expected %v, got %v", helloUserID, vp2.UserID)
	}

}

func TestShouldGetAllVisitedPathsForOneUser(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	goodbyeDate := time.Date(2018, time.November, 17, 0, 0, 0, 0, time.UTC)
	goodbyeUserID := uint32(2918592)
	helloDate := time.Date(2018, time.November, 15, 0, 0, 0, 0, time.UTC)
	helloUserID := uint32(2918592)

	sentRows := sqlmock.NewRows([]string{"path", "visit_date", "user_id"}).
		AddRow("/goodbye", goodbyeDate, goodbyeUserID).
		AddRow("/hello", helloDate, helloUserID)
	mock.ExpectQuery("[SELECT path, visit_date, user_id FROM visitedpaths WHERE user_id = $1 ORDER BY visit_date DESC]").
		WithArgs(2918592).
		WillReturnRows(sentRows)

	// run the tested function
	gotRows, err := db.GetAllVisitedPathsForUserID(uint32(2918592))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values
	if len(gotRows) != 2 {
		t.Fatalf("expected len %d, got %d", 3, len(gotRows))
	}
	vp0 := gotRows[0]
	if vp0.Path != "/goodbye" {
		t.Errorf("expected %v, got %v", "/goodbye", vp0.Path)
	}
	if vp0.Date != goodbyeDate {
		t.Errorf("expected %v, got %v", goodbyeDate, vp0.Date)
	}
	if vp0.UserID != goodbyeUserID {
		t.Errorf("expected %v, got %v", goodbyeUserID, vp0.UserID)
	}
	vp2 := gotRows[1]
	if vp2.Path != "/hello" {
		t.Errorf("expected %v, got %v", "/hello", vp2.Path)
	}
	if vp2.Date != helloDate {
		t.Errorf("expected %v, got %v", helloDate, vp2.Date)
	}
	if vp2.UserID != helloUserID {
		t.Errorf("expected %v, got %v", helloUserID, vp2.UserID)
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
	helloUserID := uint32(582)

	regexStmt := `[INSERT INTO visitedpaths(path, visit_date, user_id) VALUES (\$1, \$2, \$3)]`
	mock.ExpectPrepare(regexStmt)
	stmt := "INSERT INTO visitedpaths"
	mock.ExpectExec(stmt).
		WithArgs("hello", helloDate, helloUserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// run the tested function
	err = db.AddVisitedPath("hello", helloDate, helloUserID)
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
		Path:   "/abc",
		Date:   time.Date(2018, time.November, 17, 0, 0, 0, 0, time.UTC),
		UserID: 483,
	}

	js, err := json.Marshal(vp)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// read back in as VisitedPath struct
	vpGot := &VisitedPath{}
	err = json.Unmarshal(js, &vpGot)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// check for expected values
	if vpGot.Path != "/abc" {
		t.Errorf("expected %v, got %v", "/abc", vpGot.Path)
	}
	if vpGot.Date != vp.Date {
		t.Errorf("expected %v, got %v", vp.Date, vpGot.Date)
	}
	if vpGot.UserID != 483 {
		t.Errorf("expected %v, got %v", 483, vpGot.UserID)
	}
}

func TestVisitedPathCanUnmarshalFromJSON(t *testing.T) {
	vp := &VisitedPath{}
	js := []byte(`{"path":"/def", "date":"2018-11-17T20:43:00Z", "user_id": 5872}`)

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
	if vp.UserID != 5872 {
		t.Errorf("expected %v, got %v", 5872, vp.UserID)
	}
}
