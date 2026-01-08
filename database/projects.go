package database

import (
	"database/sql"
	"fmt"
	"strings"
)

type Project struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	SessionID   string `json:"session_id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Domain      string `json:"domain"`
	GhRepo      string `json:"gh_repo"`
	Status      string `json:"status"`
	ErrorMsg    string `json:"error_msg"`
	CreatedAt   string `json:"created_at"`
	FlyHostname string `json:"fly_hostname"`
}

func GetProjects() ([]Project, error) {
	var projects []Project
	rows, err := DB.Query(`SELECT id, name, domain, status, gh_repo
  FROM projects ORDER BY created_at DESC;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Domain, &p.Status, &p.GhRepo); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}

func UpdateProjectName(projectID, name string) error {
	_, err := DB.Exec(`
		UPDATE projects SET name = $1 WHERE id = $2;`,
		name, projectID)
	return err
}

func GetDeployedProjects() ([]Project, error) {
	var projects []Project
	rows, err := DB.Query(`SELECT id, name, domain, status, gh_repo
  FROM projects WHERE domain != '' ORDER BY created_at DESC LIMIT 20;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Domain, &p.Status, &p.GhRepo); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}

func DeleteProject(projectID string) error {
	result, err := DB.Exec(`DELETE FROM projects WHERE id = $1;`, projectID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No project found with the id %v", projectID)
	}
	return nil
}

func UpdateProjectDomain(projectID, domain string) error {
	_, err := DB.Exec(`
		UPDATE projects SET domain = $1 WHERE id = $2;`,
		domain, projectID)
	return err
}

func UpdateProjectPort(projectID string, port int) error {
	_, err := DB.Exec(`
		UPDATE projects SET port = $1 WHERE id = $2;`,
		port, projectID)
	return err
}

func UpdateGhRepo(projectID, gh_repo string) error {
	_, err := DB.Exec(`
		UPDATE projects SET gh_repo = $1 WHERE id = $2;`,
		gh_repo, projectID)
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
    gh_repo,
    error_msg,
    created_at,
    fly_hostname
    FROM projects WHERE id = $1`, id)

	if err := row.Scan(&p.ID, &p.UserID, &p.SessionID, &p.Name, &p.Slug,
		&p.Domain, &p.Status, &p.GhRepo,
		&p.ErrorMsg, &p.CreatedAt, &p.FlyHostname); err != nil {
		if err == sql.ErrNoRows {
			return p, fmt.Errorf("No project found with the id: %s", id)
		}
		return p, err
	}
	return p, nil
}

func UpdateProjectErrorMsg(projectID, msg string) error {
	_, err := DB.Exec(`
		UPDATE projects SET error_msg = $1 WHERE id = $2;`,
		msg, projectID)

	return err
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

func UpdateProjectOwner(projectID, userID string) error {
	result, err := DB.Exec(`
		UPDATE projects SET user_id = $1 WHERE id = $2;`,
		userID, projectID)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No project found with the id %v", projectID)
	}
	return err
}

func UpdateFlyHostname(projectID, flyHostname string) error {
	result, err := DB.Exec(`
		UPDATE projects SET fly_hostname = $1 WHERE id = $2;`,
		flyHostname, projectID)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No project found with the id %v", projectID)
	}
	return err
}

func CreateProject(id, userID, name, slug string) error {
	_, err := DB.Exec(`
		INSERT INTO projects (id, user_id, status, name, slug) 
    VALUES ($1, $2, $3, $4, $5)`, id, userID, "new_pending", name, slug)
	if err != nil {
		if strings.Contains(err.Error(), `pq: duplicate key value violates unique constraint "projects_slug_key"`) {
			slug = fmt.Sprintf("%s-%s", slug, id)
			_, secondTryErr := DB.Exec(`
		    INSERT INTO projects (id, user_id, status, name, slug)
        VALUES ($1, $2, $3, $4, $5)`, id, userID, "new_pending", name, slug)
			if secondTryErr != nil {
				return secondTryErr
			}
		} else {
			return err
		}
	}
	return nil
}
