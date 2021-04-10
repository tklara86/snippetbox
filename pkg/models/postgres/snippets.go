package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/tklara86/snippetbox/pkg/models"
)

// SnippetModel type wraps sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title, content string, expires time.Time) (int, error) {

	var id int
	stmt := `INSERT INTO snippets (title, content, created, expires) VALUES ( $1, $2, $3, $4) RETURNING id`

	// You can also use db.Exec which returns LastInsertId() and RowsAffected() but
	// LastInsertId() is not supported in Postgresql.
	err := m.DB.QueryRow(stmt, title, content, time.Now(), expires).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {

	// Use the QueryRow() method on the connection pool to execute our
	// SQL statement, passing in the untrusted id variable as the value for the
	// placeholder parameter. This returns a pointer to a sql.Row object which
	// holds the result from the databas
	row := m.DB.QueryRow(`SELECT * FROM snippets WHERE expires > $1 AND id = $2`, time.Now(), id)

	// Initialize a pointer to a new zeroed Snippet struct.
	s := &models.Snippet{}

	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct. Notice that the arguments
	// to row.Scan are *pointers* to the place you want to copy the data into,
	// and the number of arguments must be exactly the same as the number of
	// columns returned by your statement.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	// If everything went OK then return the Snippet object.
	return s, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil,nil
}


