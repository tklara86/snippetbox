package postgres

import (
	"database/sql"
	"errors"
	"log"
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

	// Initialize a pointer to a new zeroed Snippet struct.
	s := &models.Snippet{}

	// Use the QueryRow() method on the connection pool to execute our
	// SQL statement, passing in the untrusted id variable as the value for the
	// placeholder parameter. This returns a pointer to a sql.Row object which
	// holds the result from the database

	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct. Notice that the arguments
	// to row.Scan are *pointers* to the place you want to copy the data into,
	// and the number of arguments must be exactly the same as the number of
	// columns returned by your statement.
	//err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires
	err := m.DB.QueryRow(`SELECT * FROM snippets WHERE expires > $1 AND id = $2`, time.Now(), id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

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

	rows, err := m.DB.Query(`SELECT * FROM snippets WHERE expires > $1 ORDER BY created DESC LIMIT 10`, time.Now())

	if err != nil {
		return nil, err
	}
	//app := config.AppConfig{}

	// We defer rows.Close() to ensure the sql.Rows resultset is
	// always properly closed before the Latest() method returns. This defer
	// statement should come *after* you check for an error from the Query()
	// method. Otherwise, if Query() returns an error, you'll get a panic
	// trying to close a nil resultset.
	defer func() {
		if err = rows.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// Initialize an empty slice to hold the models.Snippets objects.
	var snippets []*models.Snippet

	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by the
	// rows.Scan() method. If iteration over all the rows completes then the
	// resultset automatically closes itself and frees-up the underlying
	// database connection.
	for rows.Next() {
		// Create a pointer to a new zeroed Snippet struct.
		s := &models.Snippet{}
		// Use rows.Scan() to copy the values from each field in the row to the
		// new Snippet object that we created. Again, the arguments to row.Scan()
		// must be pointers to the place you want to copy the data into, and the
		// number of arguments must be exactly the same was the number of
		// columns returned by your statement.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		snippets = append(snippets, s)
	}
	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the Snippets slice.
	return snippets, nil

}


