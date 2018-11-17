package models

import "time"

// VisitedPath describes an instance when this Path was requested.
type VisitedPath struct {
	Path string
	Date time.Time
}

func (db *DB) CreateTableVisitedPath() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS visitedpaths (
			id SERIAL NOT NULL PRIMARY KEY,
			path TEXT NOT NULL,
			visit_date TIMESTAMP NOT NULL
		)
	`)
	return err
}

func (db *DB) GetAllVisitedPaths() ([]*VisitedPath, error) {
	rows, err := db.sqldb.Query("SELECT path, visit_date FROM visitedpaths ORDER BY visit_date DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// or []*VisitedPath{}?
	vpaths := make([]*VisitedPath, 0)
	for rows.Next() {
		// or &VisitedPath{}?
		vp := new(VisitedPath)
		err := rows.Scan(&vp.Path, &vp.Date)
		if err != nil {
			return nil, err
		}
		vpaths = append(vpaths, vp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return vpaths, nil
}

func (db *DB) AddVisitedPath(p string, t time.Time) error {
	// move out into one-time-prepared statement?
	stmt, err := db.sqldb.Prepare("INSERT INTO visitedpaths(path, visit_date) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(p, t)
	if err != nil {
		return err
	}
	return nil
}
