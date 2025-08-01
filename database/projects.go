package database

type Project struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Status    string `json:"status"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

func UpdateProjectSessionID(projectID, sessionID string) error {
	_, err := DB.Exec(`
		UPDATE projects SET
			session_id = $1,
		WHERE id = $2;`,
		projectID, sessionID)

	return err
}

func CreateProject(id, userID, name string) error {
	_, err := DB.Exec(`
		INSERT INTO projects (id, user_id, status, name) 
    VALUES ($1, $2, $3, $4)`, id, userID, "Thinking", name)
	if err != nil {
		return err
	}
	return nil
}
