package note

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool      `json:"is_deleted"`
}

func find(id int) (*note, error) {
	db, err := sql.Open("sqlite3", "./data/lethe.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	row := db.QueryRow("SELECT * FROM notes WHERE id = ?", id)

	return instanceFromRow(row, false)
}

func all() ([]*note, error) {
	db, err := sql.Open("sqlite3", "./data/lethe.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM notes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := make([]*note, 0)
	for rows.Next() {
		note, err := instanceFromRows(rows, false)
		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}

	return notes, nil
}

func instanceFromRows(rows *sql.Rows, withDeleted bool) (*note, error) {
	var note note
	err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.UpdatedAt, &note.IsDeleted)
	if err != nil {
		return nil, err
	}

	if !withDeleted && note.IsDeleted {
		return nil, nil
	}

	return &note, nil
}

func instanceFromRow(row *sql.Row, withDeleted bool) (*note, error) {
	var note note
	err := row.Scan(&note.ID, &note.Title, &note.Content, &note.UpdatedAt, &note.IsDeleted)
	if err != nil {
		return nil, err
	}

	if !withDeleted && note.IsDeleted {
		return nil, nil
	}

	return &note, nil
}
