package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/swinslow/containerapp/api/models"
)

// ===== /admin/history route =====

func TestAdminCanGetHistory(t *testing.T) {
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

func TestCannotGetHistoryWithoutAdminUserInContext(t *testing.T) {
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

// ===== /admin/users GET route =====

func TestAdminCanGetAllUsers(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/admin/users", nil)
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
	http.HandlerFunc(env.getUsersHandler).ServeHTTP(rec, req)

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
	// read back in as slice of Users
	var users []*models.User
	err = json.Unmarshal([]byte(rec.Body.String()), &users)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// check for expected length and values
	if len(users) != 2 {
		t.Fatalf("expected len %d, got %d", 2, len(users))
	}

	user0 := users[0]
	if user0.ID != 91461 {
		t.Errorf("expected %v, got %v", 91461, user0.ID)
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

	user1 := users[1]
	if user1.ID != 914611345 {
		t.Errorf("expected %v, got %v", 914611345, user1.ID)
	}
	if user1.Email != "janedoe@example.com" {
		t.Errorf("expected %v, got %v", "janedoe@example.com", user1.Email)
	}
	if user1.Name != "Jane Doe" {
		t.Errorf("expected %v, got %v", "Jane Doe", user1.Name)
	}
	if user1.IsAdmin != true {
		t.Errorf("expected %v, got %v", true, user1.IsAdmin)
	}
}

func TestCannotGetAllUsersWithNoUserInContext(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/admin/users", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	db := &mockDB{}
	env := Env{db: db, jwtSecretKey: "keyForTesting"}
	// not adding any User to context
	http.HandlerFunc(env.getUsersHandler).ServeHTTP(rec, req)

	confirmRecWasInvalidAuth(t, rec, "Authentication header with valid Bearer token required")
}

func TestCannotGetUsersWithoutValidUserInContext(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/admin/users", nil)
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
	http.HandlerFunc(env.getUsersHandler).ServeHTTP(rec, req)

	confirmRecWasInvalidAuth(t, rec, "unknown user unknown@example.com")
}

func TestCannotGetUsersWithoutAdminUserInContext(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/admin/users", nil)
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
	http.HandlerFunc(env.getUsersHandler).ServeHTTP(rec, req)

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

// ===== /admin/users POST route =====
type mockNoUsersDB struct {
	addedUsers []*models.User
}

func (mdb *mockNoUsersDB) GetAllUsers() ([]*models.User, error) {
	users := make([]*models.User, 0)
	return users, nil
}

func (mdb *mockNoUsersDB) GetUserByID(id uint32) (*models.User, error) {
	return nil, nil
}

func (mdb *mockNoUsersDB) GetUserByEmail(email string) (*models.User, error) {
	return nil, nil
}

func (mdb *mockNoUsersDB) AddUser(id uint32, email string, name string, isAdmin bool) error {
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

func (mdb *mockNoUsersDB) GetAllVisitedPaths() ([]*models.VisitedPath, error) {
	return nil, nil
}

func (mdb *mockNoUsersDB) GetAllVisitedPathsForUserID(uint32) ([]*models.VisitedPath, error) {
	return nil, nil
}

func (mdb *mockNoUsersDB) AddVisitedPath(p string, ti time.Time, userID uint32) error {
	return nil
}
func TestCanPostNewAdminUserWithoutAuthIfNoUsers(t *testing.T) {
	rec := httptest.NewRecorder()
	body := `{"email": "steve@example.com", "name": "Steve"}`
	req, err := http.NewRequest("POST", "/admin/users", strings.NewReader(body))
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	noUsersDb := &mockNoUsersDB{}
	env := Env{db: noUsersDb, jwtSecretKey: "keyForTesting"}
	// no user in context b/c brand new
	http.HandlerFunc(env.newUserHandler).ServeHTTP(rec, req)

	// check that we got a 201 (Created)
	if 201 != rec.Code {
		t.Errorf("Expected %d, got %d", 201, rec.Code)
	}

	// check that content type was application/json
	header := rec.Result().Header
	if header.Get("Content-Type") != "application/json" {
		t.Errorf("expected %v, got %v", "application/json", header.Get("Content-Type"))
	}

	// check that the correct JSON strings were returned
	// read back in as slice of Users
	var newUser *models.User
	err = json.Unmarshal([]byte(rec.Body.String()), &newUser)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// don't check exact ID since it will be randomly generated
	// just make sure it's greater than 0
	if newUser.ID == 0 {
		t.Errorf("expected non-zero ID, got 0")
	}
	if newUser.Email != "steve@example.com" {
		t.Errorf("expected %v, got %v", "steve@example.com", newUser.Email)
	}
	if newUser.Name != "Steve" {
		t.Errorf("expected %v, got %v", "Steve", newUser.Name)
	}
	// AND this time it'll be an admin user since we're bootstrapping
	if newUser.IsAdmin != true {
		t.Errorf("expected %v, got %v", true, newUser.IsAdmin)
	}

	// and make sure that a new user was saved to database
	if len(noUsersDb.addedUsers) != 1 {
		t.Fatalf("expected 1 added user, got 0")
	}
	userCheck := noUsersDb.addedUsers[0]
	if userCheck.ID != newUser.ID {
		t.Errorf("expected same IDs, userCheck is %v, newUser is %v", userCheck.ID, newUser.ID)
	}
	if userCheck.Email != newUser.Email {
		t.Errorf("expected same email, userCheck is %v, newUser is %v", userCheck.Email, newUser.Email)
	}
	if userCheck.Name != newUser.Name {
		t.Errorf("expected same name, userCheck is %v, newUser is %v", userCheck.Name, newUser.Name)
	}
	if userCheck.IsAdmin != newUser.IsAdmin {
		t.Errorf("expected same admin status, userCheck is %v, newUser is %v", userCheck.IsAdmin, newUser.IsAdmin)
	}

}

func TestAdminCanPostNewUser(t *testing.T) {
	rec := httptest.NewRecorder()
	body := `{"email": "steve@example.com", "name": "Steve"}`
	req, err := http.NewRequest("POST", "/admin/users", strings.NewReader(body))
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
	http.HandlerFunc(env.newUserHandler).ServeHTTP(rec, req)

	// check that we got a 201 (Created)
	if 201 != rec.Code {
		t.Errorf("Expected %d, got %d", 201, rec.Code)
	}

	// check that content type was application/json
	header := rec.Result().Header
	if header.Get("Content-Type") != "application/json" {
		t.Errorf("expected %v, got %v", "application/json", header.Get("Content-Type"))
	}

	// check that the correct JSON strings were returned
	// read back in as slice of Users
	var newUser *models.User
	err = json.Unmarshal([]byte(rec.Body.String()), &newUser)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// don't check exact ID since it will be randomly generated
	// just make sure it's greater than 0
	if newUser.ID == 0 {
		t.Errorf("expected non-zero ID, got 0")
	}
	if newUser.Email != "steve@example.com" {
		t.Errorf("expected %v, got %v", "steve@example.com", newUser.Email)
	}
	if newUser.Name != "Steve" {
		t.Errorf("expected %v, got %v", "Steve", newUser.Name)
	}
	if newUser.IsAdmin != false {
		t.Errorf("expected %v, got %v", false, newUser.IsAdmin)
	}

	// and make sure that a new user was saved to database
	if len(db.addedUsers) != 1 {
		t.Fatalf("expected 1 added user, got 0")
	}
	userCheck := db.addedUsers[0]
	if userCheck.ID != newUser.ID {
		t.Errorf("expected same IDs, userCheck is %v, newUser is %v", userCheck.ID, newUser.ID)
	}
	if userCheck.Email != newUser.Email {
		t.Errorf("expected same email, userCheck is %v, newUser is %v", userCheck.Email, newUser.Email)
	}
	if userCheck.Name != newUser.Name {
		t.Errorf("expected same name, userCheck is %v, newUser is %v", userCheck.Name, newUser.Name)
	}
	if userCheck.IsAdmin != newUser.IsAdmin {
		t.Errorf("expected same admin status, userCheck is %v, newUser is %v", userCheck.IsAdmin, newUser.IsAdmin)
	}

}

func TestAdminCannotPostNewUserWithExistingEmail(t *testing.T) {
	rec := httptest.NewRecorder()
	body := `{"email": "johndoe@example.com", "name": "oops John Doe Redux"}`
	req, err := http.NewRequest("POST", "/admin/users", strings.NewReader(body))
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
	http.HandlerFunc(env.newUserHandler).ServeHTTP(rec, req)

	// check that we got a 400 (Bad Request)
	if 400 != rec.Code {
		t.Errorf("Expected %d, got %d", 400, rec.Code)
	}

	// check that content type was application/json
	header := rec.Result().Header
	if header.Get("Content-Type") != "application/json" {
		t.Errorf("expected %v, got %v", "application/json", header.Get("Content-Type"))
	}

	// check that the right "error" JSON string was returned
	wantString := `{"error": "user with email johndoe@example.com already exists"}`
	if rec.Body.String() != wantString {
		t.Fatalf("expected %s, got %s", wantString, rec.Body.String())
	}
}

func TestCannotPostNewUserWithNoUserInContext(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/admin/users", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	db := &mockDB{}
	env := Env{db: db, jwtSecretKey: "keyForTesting"}
	// not adding any User to context
	http.HandlerFunc(env.newUserHandler).ServeHTTP(rec, req)

	confirmRecWasInvalidAuth(t, rec, "Authentication header with valid Bearer token required")
}

func TestCannotPostNewUserWithoutValidUserInContext(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/admin/users", nil)
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
	http.HandlerFunc(env.newUserHandler).ServeHTTP(rec, req)

	confirmRecWasInvalidAuth(t, rec, "unknown user unknown@example.com")
}

func TestCannotPostNewUserWithoutAdminUserInContext(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/admin/users", nil)
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
	http.HandlerFunc(env.newUserHandler).ServeHTTP(rec, req)

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
