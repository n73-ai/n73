package database

type Message struct {
	ID         string `json:"id"`
	ProjectID  string `json:"project_id"`
	MetadataID string `json:"metadata_id"`
	Role       string `json:"role"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
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
