package database

import (
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
		VALUES (?, ?, ?, ?, ?)`,
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
