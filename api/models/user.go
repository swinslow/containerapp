package models

import "fmt"

// User describes a registered user of the platform.
type User struct {
	ID      uint32 `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"is_admin"`
}

// CreateTableUsers creates the users table if it does not already exist.
func (db *DB) CreateTableUsers() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER NOT NULL PRIMARY KEY,
			email TEXT NOT NULL,
			name TEXT NOT NULL,
			is_admin BOOLEAN NOT NULL
		)
	`)
	return err
}

// GetAllUsers returns a slice with all registered users.
func (db *DB) GetAllUsers() ([]*User, error) {
	rows, err := db.sqldb.Query("SELECT id, email, name, is_admin FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// or []*User{}?
	users := make([]*User, 0)
	for rows.Next() {
		// or &User{}?
		user := new(User)
		err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.IsAdmin)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// AddUser adds a user to the database.
// Due to PostgreSQL limits on integer size, id must be less than 2147483647.
// It should typically be created via math/rand's Int31() function and then
// cast to uint32.
func (db *DB) AddUser(id uint32, email string, name string, isAdmin bool) error {
	var maxUserID uint32
	maxUserID = 2147483647

	if id > maxUserID {
		return fmt.Errorf("User id cannot be greater than %d; received %d", maxUserID, id)
	}

	// move out into one-time-prepared statement?
	stmt, err := db.sqldb.Prepare("INSERT INTO users(id, email, name, is_admin) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id, email, name, isAdmin)
	if err != nil {
		return err
	}
	return nil
}
