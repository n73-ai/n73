package database

import (
	"fmt"
)

type Message struct {
	ID           string  `json:"id"`
	ProjectID    string  `json:"project_id"`
	Role         string  `json:"role"`
	Content      string  `json:"content"`
	Model        string  `json:"model"`
	Duration     int     `json:"duration"`
	IsError      bool    `json:"is_error"`
	TotalCostUsd float64 `json:"total_cost_usd"`
	CreatedAt    string  `json:"created_at"`
}

func GetMessagesByProjectID(projectID string) ([]Message, error) {
	var messages []Message
	rows, err := DB.Query(`SELECT id, role, content, duration, total_cost_usd 
  FROM messages WHERE project_id = $1;`, projectID)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.Role, &m.Content, &m.Duration, &m.TotalCostUsd); err != nil {
			return nil, fmt.Errorf("error: %v", err)
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	return messages, nil
}

func CreateMessage(id, projectID, role, content, model string, duration int, isError bool, totalCostUsd float64) error {
	_, err := DB.Exec(`
		INSERT INTO messages 
    (id, project_id, role, content, model, duration, is_error, total_cost_usd) 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		id, projectID, role, content, model, duration, isError, totalCostUsd)
	if err != nil {
		return err
	}
	return nil
}
