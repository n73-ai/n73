package database

import (
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

func CreateLog(errorScope, entityID, message string) {
	id := uuid.NewString()
	_, err := DB.Exec(`
		INSERT INTO logs (id, error_scope, entity_id, message)
    VALUES ($1, $2, $3, $4)`, id, errorScope, entityID, message)
	if err != nil {
		log.Errorf("CreateLog() error: %v", err)
	}
}
