package database

import (
	"database/sql"
	"fmt"
)

type Project struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	SessionID      string `json:"session_id"`
	Status         string `json:"status"`
	CFProjectReady bool   `json:"cf_project_ready"`
	Name           string `json:"name"`
	Domain         string `json:"domain"`
	CreatedAt      string `json:"created_at"`
}

func GetProjectsByUserID(userID string) ([]Project, error) {
	var projects []Project
	rows, err := DB.Query(`SELECT id, name, domain
  FROM projects WHERE user_id = $1 ORDER BY created_at DESC;`, userID)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Domain); err != nil {
			return nil, fmt.Errorf("error: %v", err)
		}
		projects = append(projects, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	return projects, nil
}

func GetProjectByID(id string) (Project, error) {
	var p Project
	row := DB.QueryRow(`SELECT id, user_id, session_id, status, cf_project_ready,
        name, domain, created_at FROM projects WHERE id = $1`, id)

	if err := row.Scan(&p.ID, &p.UserID, &p.SessionID, &p.Status,
		&p.CFProjectReady, &p.Name, &p.Domain, &p.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return p, fmt.Errorf("No project found with the id: %s", id)
		}
		return p, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	return p, nil
}

func UpdateProjectStatus(projectID, status string) error {
	_, err := DB.Exec(`
		UPDATE projects SET status = $1 WHERE id = $2;`,
		status, projectID)

	return err
}

func UpdateProjectSessionID(projectID, sessionID string) error {
	_, err := DB.Exec(`
		UPDATE projects SET session_id = $1 WHERE id = $2;`,
		sessionID, projectID)

	return err
}

func CreateProject(id, userID, name string) error {
	_, err := DB.Exec(`
		INSERT INTO projects (id, user_id, status, name) 
    VALUES ($1, $2, $3, $4)`, id, userID, "Building", name)
	if err != nil {
		return err
	}
	return nil
}
