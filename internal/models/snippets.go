package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID 				int
	Title 		string
	Content 	string
	Created 	time.Time
	Expires 	time.Time
}

// SnippetModel type that wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	sqlQuery := `
		INSERT INTO snippets (title, content, created, expires)
		VALUES($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + ($3 * INTERVAL '1 day'))
		RETURNING id;
	`

	var id int64
	// Scan helps us get the column 'id' from the row we get from QueryRow.
	err := m.DB.QueryRow(sqlQuery, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}

	// Convert int64 to a int32.
	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (Snippet, error) {
	sqlQuery := `
		SELECT id, title, content, created, expires FROM snippets
		WHERE expires > CURRENT_TIMESTAMP and id = $1;
	`
	row := m.DB.QueryRow(sqlQuery, id)

	var s Snippet
	// It will automatically convert the raw output from the SQL database
	// to the required native Go types.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return s, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]Snippet, error) {
	sqlQuery := `
		SELECT id, title, content, created, expires FROM snippets
		WHERE expires > CURRENT_TIMESTAMP
		ORDER BY id DESC
		LIMIT 10;
	`
	rows, err := m.DB.Query(sqlQuery)
	if err != nil {
		return nil, err
	}	

	defer rows.Close()
	var snippets []Snippet

	for rows.Next() {
		var s Snippet
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}