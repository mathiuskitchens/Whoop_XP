package whoop

import (
	"database/sql"
	"fmt"
)

type Metric struct {
	ID            int    `json:"id"`
	Date          string `json:"date" binding:"required"`
	SleepScore    int    `json:"sleep_score" binding:"required"`
	RecoveryScore int    `json:"recovery_score" binding:"required"`
	StrainScore   int    `json:"strain_score" binding:"required"`
}

func GetAllMetrics(db *sql.DB) ([]Metric, error) {
	rows, err := db.Query(`
		SELECT id, date, sleep_score, recovery_score, strain_score
		FROM whoop_metrics
		ORDER BY date DESC
		`)

	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}

	defer rows.Close()

	var metrics []Metric

	for rows.Next() {
		var m Metric
		err := rows.Scan(&m.ID, &m.Date, &m.SleepScore, &m.RecoveryScore, &m.StrainScore)
		if err != nil {
			return nil, fmt.Errorf("row scan failed: %v", err)
		}
		metrics = append(metrics, m)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed; %v", err)
	}

	return metrics, nil
}

func InsertMetric(db *sql.DB, m Metric) error {
	query := `
		INSERT INTO whoop_metrics (date, sleep_score, recovery_score, strain_score)
		VALUES (?, ?, ?, ?)
	`

	_, err := db.Exec(query, m.Date, m.SleepScore, m.RecoveryScore, m.StrainScore)
	if err != nil {
		return fmt.Errorf("insert metric mailed: %v", err)
	}

	fmt.Println("Metric inserted!")
	return nil
}
