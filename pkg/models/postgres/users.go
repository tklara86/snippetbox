package postgres

import (
	"database/sql"
	"github.com/tklara86/snippetbox/pkg/models"
)

type UserModel struct {
	DB    *sql.DB
}

// Insert add a new record to users table
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// Authenticate verifies whether a user exists with the provided email address and password. This will return
// relevant User ID if they do.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0,nil
}

// Get fetches details for a specific user based on their user ID
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}

