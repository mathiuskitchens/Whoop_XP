package spiritual

import (
	"database/sql"
	"fmt"
)

type Discipline struct {
	ID         int    `json:"id"`
	Date       string `json:"date" binding:"required"`
	Prayer     bool   `json:"prayer"`
	Scripture  bool   `json:"scripture"`
	Meditation bool   `json:"meditation"`
	Journaling bool   `json:"journaling"`
	Gratitude  bool   `json:"gratitude"`
}

// InsertDiscipline adds a new daily record

func InsertDiscipline(db *sql.DB, d Discipline) error {
	query := `
	INSERT INTO spiritual_disciplines (date, prayer, scripture, meditation, journaling, gratitude)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(query, d.Date, d.Prayer, d.Scripture, d.Meditation, d.Journaling, d.Gratitude)
	if err != nil {
		return fmt.Errorf("insert discipline failed: %v", err)
	}
	return nil
}

func GetAllDisciplines(db *sql.DB) ([]Discipline, error) {
	rows, err := db.Query(`
	SELECT * 
	FROM spiritual_disciplines
	ORDER BY date DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("get all disciplines failed: %v", err)
	}
	defer rows.Close()

	var disciplines []Discipline
	for rows.Next() {
		var d Discipline
		err := rows.Scan(&d.ID, &d.Date, &d.Prayer, &d.Scripture, &d.Meditation, &d.Journaling, &d.Gratitude)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %v", err)
		}
		disciplines = append(disciplines, d)
	}
	return disciplines, nil
}
