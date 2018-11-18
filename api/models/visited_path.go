package models

import (
	"encoding/json"
	"time"
)

// VisitedPath describes an instance when this Path was requested.
type VisitedPath struct {
	Path   string
	Date   time.Time
	UserID uint32
}

func (vp *VisitedPath) MarshalJSON() ([]byte, error) {
	fmtVp := struct {
		Path   string `json:"path"`
		Date   string `json:"date"`
		UserID uint32 `json:"user_id"`
	}{
		Path:   vp.Path,
		Date:   vp.Date.Format(time.RFC3339),
		UserID: vp.UserID,
	}

	return json.Marshal(fmtVp)
}

func (vp *VisitedPath) UnmarshalJSON(js []byte) error {
	var strs map[string]interface{}

	err := json.Unmarshal(js, &strs)
	if err != nil {
		return err
	}

	for k, v := range strs {
		switch k {
		case "path":
			vp.Path = v.(string)
		case "date":
			ti, err := time.Parse(time.RFC3339, v.(string))
			if err != nil {
				return err
			}
			vp.Date = ti
		case "user_id":
			vp.UserID = uint32(v.(float64))
		}
	}

	return nil
}

func (db *DB) CreateTableVisitedPath() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS visitedpaths (
			id SERIAL NOT NULL PRIMARY KEY,
			path TEXT NOT NULL,
			visit_date TIMESTAMP NOT NULL,
			user_id INTEGER NOT NULL
		)
	`)
	return err
}

func (db *DB) GetAllVisitedPaths() ([]*VisitedPath, error) {
	rows, err := db.sqldb.Query("SELECT path, visit_date, user_id FROM visitedpaths ORDER BY visit_date DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// or []*VisitedPath{}?
	vpaths := make([]*VisitedPath, 0)
	for rows.Next() {
		// or &VisitedPath{}?
		vp := new(VisitedPath)
		err := rows.Scan(&vp.Path, &vp.Date, &vp.UserID)
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

func (db *DB) GetAllVisitedPathsForUserID(user_id uint32) ([]*VisitedPath, error) {
	rows, err := db.sqldb.Query("SELECT path, visit_date, user_id FROM visitedpaths WHERE user_id = $1 ORDER BY visit_date DESC", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// or []*VisitedPath{}?
	vpaths := make([]*VisitedPath, 0)
	for rows.Next() {
		// or &VisitedPath{}?
		vp := new(VisitedPath)
		err := rows.Scan(&vp.Path, &vp.Date, &vp.UserID)
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

func (db *DB) AddVisitedPath(p string, t time.Time, user_id uint32) error {
	// move out into one-time-prepared statement?
	stmt, err := db.sqldb.Prepare("INSERT INTO visitedpaths(path, visit_date, user_id) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(p, t, user_id)
	if err != nil {
		return err
	}
	return nil
}
