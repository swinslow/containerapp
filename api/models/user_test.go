package models

import (
	"encoding/json"
	"strings"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestShouldGetAllUsers(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	sentRows := sqlmock.NewRows([]string{"id", "email", "name", "is_admin"}).
		AddRow(410952, "johndoe@example.com", "John Doe", false).
		AddRow(8103918, "janedoe@example.com", "Jane Doe", true)
	mock.ExpectQuery("SELECT id, email, name, is_admin FROM users ORDER BY id").WillReturnRows(sentRows)

	// run the tested function
	gotRows, err := db.GetAllUsers()
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
		t.Fatalf("expected len %d, got %d", 2, len(gotRows))
	}
	user0 := gotRows[0]
	if user0.ID != 410952 {
		t.Errorf("expected %v, got %v", 410952, user0.ID)
	}
	if user0.Email != "johndoe@example.com" {
		t.Errorf("expected %v, got %v", "johndoe@example.com", user0.Email)
	}
	if user0.Name != "John Doe" {
		t.Errorf("expected %v, got %v", "John Doe", user0.Name)
	}
	if user0.IsAdmin != false {
		t.Errorf("expected %v, got %v", false, user0.IsAdmin)
	}

}

func TestShouldGetUserByID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	sentRows := sqlmock.NewRows([]string{"id", "email", "name", "is_admin"}).
		AddRow(8103918, "janedoe@example.com", "Jane Doe", true)
	mock.ExpectQuery(`[SELECT id, email, name, is_admin FROM users WHERE id = \$1]`).
		WithArgs(8103918).
		WillReturnRows(sentRows)

	// run the tested function
	user, err := db.GetUserByID(8103918)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values
	if user.ID != 8103918 {
		t.Errorf("expected %v, got %v", 8103918, user.ID)
	}
	if user.Email != "janedoe@example.com" {
		t.Errorf("expected %v, got %v", "janedoe@example.com", user.Email)
	}
	if user.Name != "Jane Doe" {
		t.Errorf("expected %v, got %v", "Jane Doe", user.Name)
	}
	if user.IsAdmin != true {
		t.Errorf("expected %v, got %v", true, user.IsAdmin)
	}

}

func TestShouldGetUserByEmail(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	sentRows := sqlmock.NewRows([]string{"id", "email", "name", "is_admin"}).
		AddRow(8103918, "janedoe@example.com", "Jane Doe", true)
	mock.ExpectQuery(`[SELECT id, email, name, is_admin FROM users WHERE email = \$1]`).
		WithArgs("janedoe@example.com").
		WillReturnRows(sentRows)

	// run the tested function
	user, err := db.GetUserByEmail("janedoe@example.com")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values
	if user.ID != 8103918 {
		t.Errorf("expected %v, got %v", 8103918, user.ID)
	}
	if user.Email != "janedoe@example.com" {
		t.Errorf("expected %v, got %v", "janedoe@example.com", user.Email)
	}
	if user.Name != "Jane Doe" {
		t.Errorf("expected %v, got %v", "Jane Doe", user.Name)
	}
	if user.IsAdmin != true {
		t.Errorf("expected %v, got %v", true, user.IsAdmin)
	}

}

func TestShouldAddUser(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[INSERT INTO users(id, email, name, is_admin) VALUES (\$1, \$2, \$3, \$4)]`
	mock.ExpectPrepare(regexStmt)
	stmt := "INSERT INTO users"
	mock.ExpectExec(stmt).
		WithArgs(192304, "johndoe@example.com", "John Doe", false).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// run the tested function
	err = db.AddUser(192304, "johndoe@example.com", "John Doe", false)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldNotAddUserWithGreaterThanMaxID(t *testing.T) {
	// set up mock
	sqldb, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	// run the tested function
	err = db.AddUser(2147483648, "oops@example.com", "OOPS", false)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
	// a non-nil error that's related to sqlmock errors is still wrong
	if err != nil && strings.Contains(err.Error(), "all expectations were already fulfilled, call to Prepare") {
		t.Fatalf("didn't expect sqlmock error: %v", err)
	}
}

// ===== JSON marshalling and unmarshalling =====
func TestCanMarshalAdminUserToJSON(t *testing.T) {
	user := &User{
		ID:      85010942,
		Email:   "janedoe@example.com",
		Name:    "Jane Doe",
		IsAdmin: true,
	}

	js, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// read back in as User struct
	userGot := &User{}
	err = json.Unmarshal(js, userGot)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// check for expected length and values
	if user.ID != userGot.ID {
		t.Errorf("expected %v, got %v", user.ID, userGot.ID)
	}
	if user.Email != userGot.Email {
		t.Errorf("expected %v, got %v", user.Email, userGot.Email)
	}
	if user.Name != userGot.Name {
		t.Errorf("expected %v, got %v", user.Name, userGot.Name)
	}
	if user.IsAdmin != userGot.IsAdmin {
		t.Errorf("expected %v, got %v", user.IsAdmin, userGot.IsAdmin)
	}
}

func TestCanMarshalNonAdminUserToJSON(t *testing.T) {
	user := &User{
		ID:      16923941,
		Email:   "johndoe@example.com",
		Name:    "John Doe",
		IsAdmin: false,
	}

	js, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// read back in as User struct
	userGot := &User{}
	err = json.Unmarshal(js, userGot)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// check for expected length and values
	if user.ID != userGot.ID {
		t.Errorf("expected %v, got %v", user.ID, userGot.ID)
	}
	if user.Email != userGot.Email {
		t.Errorf("expected %v, got %v", user.Email, userGot.Email)
	}
	if user.Name != userGot.Name {
		t.Errorf("expected %v, got %v", user.Name, userGot.Name)
	}
	if user.IsAdmin != userGot.IsAdmin {
		t.Errorf("expected %v, got %v", user.IsAdmin, userGot.IsAdmin)
	}
}

func TestCanUnmarshalAdminUserFromJSON(t *testing.T) {
	user := &User{}
	js := []byte(`{"id":1920, "name":"Jane Doe", "email":"janedoe@example.com", "is_admin":true}`)

	err := json.Unmarshal(js, user)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// check values
	if user.ID != 1920 {
		t.Errorf("expected %v, got %v", 1920, user.ID)
	}
	if user.Email != "janedoe@example.com" {
		t.Errorf("expected %v, got %v", "janedoe@example.com", user.Email)
	}
	if user.Name != "Jane Doe" {
		t.Errorf("expected %v, got %v", "janedoe@example.com", user.Name)
	}
	if user.IsAdmin != true {
		t.Errorf("expected %v, got %v", true, user.IsAdmin)
	}
}

func TestCanUnmarshalNonAdminUserFromJSON(t *testing.T) {
	user := &User{}
	js := []byte(`{"id":92841, "name":"John Doe", "email":"johndoe@example.com", "is_admin":false}`)

	err := json.Unmarshal(js, user)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// check values
	if user.ID != 92841 {
		t.Errorf("expected %v, got %v", 92841, user.ID)
	}
	if user.Email != "johndoe@example.com" {
		t.Errorf("expected %v, got %v", "johndoe@example.com", user.Email)
	}
	if user.Name != "John Doe" {
		t.Errorf("expected %v, got %v", "johndoe@example.com", user.Name)
	}
	if user.IsAdmin != false {
		t.Errorf("expected %v, got %v", false, user.IsAdmin)
	}
}

func TestCannotUnmarshalUserWithNegativeIDFromJSON(t *testing.T) {
	user := &User{}
	js := []byte(`{"id":-92841, "name":"OOPS", "email":"oops@example.com", "is_admin":false}`)

	err := json.Unmarshal(js, user)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
