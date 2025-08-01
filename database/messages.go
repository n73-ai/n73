package database

import (
	"fmt"
)

type Message struct {
	ID         string `json:"id"`
	ProjectID  string `json:"project_id"`
	MetadataID string `json:"metadata_id"`
	Role       string `json:"role"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}

func GetMessagesByProjectID(projectID string) ([]Message, error) {
	var messages []Message
	rows, err := DB.Query(`SELECT id, metadata_id, role, content FROM messages WHERE project_id = $1;`, projectID)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.MetadataID, &m.Role, &m.Content); err != nil {
			return nil, fmt.Errorf("error: %v", err)
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	return messages, nil
}

func CreateMessage(id, projectID, metadatID, role, content string) error {
	_, err := DB.Exec(`
		INSERT INTO messages (id, project_id, metadata_id, role, content) 
    VALUES ($1, $2, $3, $4, $5)`, id, projectID, metadatID, role, content)
	if err != nil {
		return err
	}
	return nil
}
