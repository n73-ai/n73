package database

type Metadata struct {
	ID           string  `json:"id"`
	SessionID    string  `json:"session_id"`
	Model        string  `json:"model"`
	Duration     int     `json:"duration"`
	IsError      bool    `json:"is_error"`
	TotalCostUsd float64 `json:"total_cost_usd"`
	CreatedAt    string  `json:"created_at"`
}

func UpdateMetadata(sessionID string, duration int, totalCostUSD float64, id string, isError bool) error {
	_, err := DB.Exec(`
		UPDATE metadata SET
			session_id = $1,
			duration = $2,
			is_error = $3,
			total_cost_usd = $4,
		WHERE id = $5;`,
		sessionID, duration, isError, totalCostUSD, id)

	return err
}

func CreateMetadata(id, model string) error {
	_, err := DB.Exec(`
		INSERT INTO metadata (id, model) 
   VALUES ($1, $2);`, id, model)
	if err != nil {
		return err
	}
	return nil
}
