/*
  - project.status:
    building
    deploying
    deployed
    error
  - project.stage:
    0 none
    1 has github repository
    2 has code on github
    3 has cloudflare page
    4 has code on cloudflare
  - project.error_stage:
    0 error creating github remote repository
    1 error pushing code to github
    2 error creating cloudflare page
    3 error pushing code to cloudflare
*/
package database

import (
	"database/sql"
	"fmt"
	"strings"
)

type Project struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	SessionID  string `json:"session_id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	Domain     string `json:"domain"`
	Status     string `json:"status"`
	Stage      string `json:"stage"`
	ErrorStage string `json:"error_stage"`
	Port       int    `json:"port"`
	CreatedAt  string `json:"created_at"`
}

func UpdateProjectErrorStage(projectID, errorStage string) error {
	_, err := DB.Exec(`
		UPDATE projects SET error_stage = $1 WHERE id = $2;`,
		errorStage, projectID)
	return err
}

func UpdateProjectStage(projectID, stage string) error {
	_, err := DB.Exec(`
		UPDATE projects SET stage = $1 WHERE id = $2;`,
		stage, projectID)
	return err
}

func UpdateProjectDomain(projectID, domain string) error {
	_, err := DB.Exec(`
		UPDATE projects SET domain = $1 WHERE id = $2;`,
		domain, projectID)
	return err
}

func GetProjectsByUserID(userID string) ([]Project, error) {
	var projects []Project
	rows, err := DB.Query(`SELECT id, name, domain, status
  FROM projects WHERE user_id = $1 ORDER BY created_at DESC;`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Domain, &p.Status); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}

func GetProjectByID(id string) (Project, error) {
	var p Project
	row := DB.QueryRow(`SELECT 
    id, 
    user_id, 
    session_id, 
    name, 
    slug,
    domain, 
    status, 
    stage,
    error_stage,
    port,
    created_at 
    FROM projects WHERE id = $1`, id)

	if err := row.Scan(&p.ID, &p.UserID, &p.SessionID, &p.Name, &p.Slug, &p.Domain, &p.Status,
		&p.Stage, &p.ErrorStage, &p.Port, &p.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return p, fmt.Errorf("No project found with the id: %s", id)
		}
		return p, err
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

func CreateProject(id, userID, name string, port int) error {
	// check requirments of cf pages to be a correct slug
	slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	_, err := DB.Exec(`
		INSERT INTO projects (id, user_id, status, name, slug, port) 
    VALUES ($1, $2, $3, $4, $5, $6)`, id, userID, "Building", name, slug, port)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: project.slug") {
			slug = fmt.Sprintf("%s-%s", slug, id)
			_, secondTryErr := DB.Exec(`
		    INSERT INTO projects (id, user_id, status, name) 
        VALUES ($1, $2, $3, $4)`, id, userID, "Building", name)
			if secondTryErr != nil {
				return secondTryErr
			}
		} else {
			return err
		}
	}
	return nil
}
