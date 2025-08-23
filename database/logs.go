package database

import (
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

type Log struct {
	ID         string `json:"id"`
	ErrorScope string `json:"error_scope"`
	EntityID   string `json:"entity_id"`
	Message    string `json:"message"`
	CreatedAt  string `json:"created_at"`
}

func GetLogs() ([]Log, error) {
	var logs []Log
	rows, err := DB.Query(`SELECT message, entity_id
  FROM logs ORDER BY created_at ASC;`)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var l Log
		if err := rows.Scan(&l.Message, &l.EntityID); err != nil {
			return nil, fmt.Errorf("error: %v", err)
		}
		logs = append(logs, l)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	return logs, nil
}

func CreateLog(errorScope, entityID, message string) {
	id := uuid.NewString()
	_, err := DB.Exec(`
		INSERT INTO logs (id, error_scope, entity_id, message)
    VALUES ($1, $2, $3, $4)`, id, errorScope, entityID, message)
	if err != nil {
		log.Errorf("CreateLog() error: %v", err)
	}
}
