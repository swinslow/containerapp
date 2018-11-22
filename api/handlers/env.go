package handlers

import (
	"fmt"
	"os"

	"github.com/swinslow/containerapp/api/models"
)

// Env is the environment for the web handlers.
type Env struct {
	db           models.Datastore
	jwtSecretKey string
}

// SetupEnv sets up systems (such as the data store) and variables
// (such as the JWT signing key) that are used across web requests.
func SetupEnv() (*Env, error) {
	// set up datastore
	db, err := models.NewDB("host=db sslmode=disable dbname=dev user=postgres-dev")
	if err != nil {
		return nil, err
	}

	err = db.InitDBTables()
	if err != nil {
		return nil, err
	}

	// set up JWT secret key (from environment)
	JWTSECRETKEY := os.Getenv("JWTSECRETKEY")
	if JWTSECRETKEY == "" {
		return nil, fmt.Errorf("No secret key found; set environment variable JWTSECRETKEY before starting")
	}

	env := &Env{
		db:           db,
		jwtSecretKey: JWTSECRETKEY,
	}
	return env, nil
}
