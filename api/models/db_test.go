package models

import "testing"

// FIXME not sure whether this works if user doesn't actually have
// FIXME postgres installed locally...
func TestCanCreateNewDBObject(t *testing.T) {
	db, err := NewDB("sslmode=disable")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if db == nil {
		t.Fatalf("expected non-nil db, got nil")
	}
}

func TestCannotCreateNewDBIfError(t *testing.T) {
	db, err := NewDB("dbname=FAILThisIsNotADB sslmode=disable")
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	if db != nil {
		t.Fatalf("expected nil db, got %v", db)
	}
}
