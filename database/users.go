package database

import (
	"database/sql"
	"fmt"
	"strings"
)

type User struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	Active    bool    `json:"active"`
	Role      string  `json:"role"`
	Balance   float64 `json:"balance"`
	CreatedAt string  `json:"created_at"`
}

func GetUserByEmail(email string) (User, error) {
	var u User
	row := DB.QueryRow(`SELECT id, email, active, role, balance, created_at 
  FROM users WHERE email = $1`, email)
	if err := row.Scan(&u.ID, &u.Email, &u.Active, &u.Role,
		&u.Balance, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return u, fmt.Errorf("No user found with the email: %s.", email)
		}
		return u, fmt.Errorf("An unexpected error occurred: %v.", err)
	}
	return u, nil
}

func UserExists(email string) (bool, error) {
	var exists bool
	err := DB.QueryRow(`SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`,
		email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func CreateUser(email, id string) error {
	_, err := DB.Exec(`
		INSERT INTO users
		(id, email, active, role, balance) 
		VALUES ($1, $2, $3, $4, $5)`,
		id, email, true, "user", 0.0)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
			return fmt.Errorf("An account with the email %s already exists. Please log in.", email)
		} else {
			return err
		}
	}

	return nil
}
