package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"ptodd.org/snippetbox/pkg/models"
)

// UserModel wraps a database connection pool
type UserModel struct {
	DB *sql.DB
}

// Insert adds a new record to the users table
func (m *UserModel) Insert(name, email, password string) error {

	// Hash the plain-text password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	// Insert SQL to add a row into the users table
	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES(?, ?, ?, UTC_TIMESTAMP())`

	// Execute the insert
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		// If this returns an error, we use the errors.As() function to check
		// whether the error has the type *mysql.MySQLError. If it does, the
		// error will be assigned to the mySQLError variable. We can then check
		// whether or not the error relates to our users_uc_email key by
		// checking the contents of the message string. If it does, we return
		// an ErrDuplicateEmail error.
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
			// Other specific SQL errors should be handled here
		}
		return err
	}

	return nil
}

// Authenticate verifies whether a user exists with
// the provided email address and password. This will return the relevant
// user ID if they do.
func (m *UserModel) Authenticate(email, password string) (int, error) {

	// Retrieve the id and hashed password associated with the given email. If no
	// matching email exists, or the user is not active, we return the
	// ErrInvalidCredentials error.
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, hashed_password FROM users WHERE email = ? AND active = TRUE"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil && errors.Is(err, sql.ErrNoRows) { // user not found
		return 0, models.ErrInvalidCredentials
	}
	if err != nil { // all other erros
		return 0, err
	}

	// Check whether the hashed password and plain-text password provided match.
	// If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) { // password does not match
		return 0, models.ErrInvalidCredentials
	}
	if err != nil { // all other errors
		return 0, err
	}

	return id, nil
}

// Get fetch details for a specific user based on their user ID.
func (m *UserModel) Get(id int) (*models.User, error) {

	// Instantiate new user model
	u := &models.User{}

	// Select SQL to retrieve a specific user from the database using their user ID
	stmt := `SELECT id, name, email, created, active FROM users WHERE id = ?`

	// Query the database and handle any errors
	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Name, &u.Email, &u.Created, &u.Active)
	if err != nil && errors.Is(err, sql.ErrNoRows) { // No record found
		return nil, models.ErrNoRecord
	}
	if err != nil { // All other unhandled errors
		return nil, err
	}

	// Handle response errors

	return u, nil
}
