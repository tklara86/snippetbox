package postgres

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/tklara86/snippetbox/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserModel struct {
	DB   *sql.DB
}

// Insert add a new record to users table
func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password),12)
	if err != nil {
		return  err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES ($1, $2, $3, $4)`

	_,err = m.DB.Exec(stmt, name, email, string(hashedPassword), time.Now())

	if err != nil {
		var postgreSQLError *pq.Error
		if errors.As(err, &postgreSQLError) {
				return  models.ErrDuplicateEmail
		}

	}
	return  nil
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

