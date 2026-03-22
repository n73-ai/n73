package database

import "time"

type FlyioIncident struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Impact    string    `json:"impact"`
	Resolved  bool      `json:"resolved"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func UpsertFlyioIncident(id, name, status, impact string, resolved bool) error {
	_, err := DB.Exec(`
		INSERT INTO flyio_incidents (id, name, status, impact, resolved, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (id) DO UPDATE SET name = $2, status = $3, impact = $4, resolved = $5, updated_at = NOW();`,
		id, name, status, impact, resolved)
	return err
}

func GetActiveFlyioIncidents() ([]FlyioIncident, error) {
	rows, err := DB.Query(`
		SELECT id, name, status, impact, resolved, created_at, updated_at
		FROM flyio_incidents
		WHERE resolved = false
		ORDER BY created_at DESC;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []FlyioIncident
	for rows.Next() {
		var i FlyioIncident
		err := rows.Scan(&i.ID, &i.Name, &i.Status, &i.Impact, &i.Resolved, &i.CreatedAt, &i.UpdatedAt)
		if err != nil {
			return nil, err
		}
		incidents = append(incidents, i)
	}
	return incidents, nil
}

func HasActiveMajorFlyioIncident() (bool, error) {
	var count int
	err := DB.QueryRow(`
		SELECT COUNT(*) FROM flyio_incidents
		WHERE resolved = false AND impact IN ('major', 'critical');`).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
