package db

import (
	"database/sql"
	"time"

	"github.com/turbo-d/overuse/models"
)

func (db Database) GetAllActivities() (*models.ActivityList, error) {
	list := &models.ActivityList{}
	rows, err := db.Conn.Query("SELECT * FROM activities ORDER BY ID DESC")
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var activity models.Activity
		err := rows.Scan(&activity.ID, &activity.Name, &activity.Description, &activity.Units, &activity.IsDeleted, &activity.IsArchived, &activity.CreatedAt)
		if err != nil {
			return list, err
		}
		list.Activities = append(list.Activities, activity)
	}
	return list, nil
}

func (db Database) AddActivity(activity *models.Activity) error {
	var id int
	var createdAt time.Time
	query := `INSERT INTO activities (name, description, units) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := db.Conn.QueryRow(query, activity.Name, activity.Description, activity.Units).Scan(&id, &createdAt)
	if err != nil {
		return err
	}
	activity.ID = id
	activity.CreatedAt = createdAt
	return nil
}

func (db Database) GetActivityById(activityID int) (models.Activity, error) {
	activity := models.Activity{}
	query := `SELECT * FROM activities WHERE id = $1;`
	row := db.Conn.QueryRow(query, activityID)
	switch err := row.Scan(&activity.ID, &activity.Name, &activity.Description, &activity.Units, &activity.IsDeleted, &activity.IsArchived, &activity.CreatedAt); err {
	case sql.ErrNoRows:
		return activity, ErrNoMatch
	default:
		return activity, err
	}
}

func (db Database) DeleteActivity(activityID int) error {
	// TODO: Only mark as deleted? Or do I do that somewhere else?
	query := `DELETE FROM activities WHERE id = $1;`
	_, err := db.Conn.Exec(query, activityID)
	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}

func (db Database) UpdateActivity(activityID int, activityData models.Activity) (models.Activity, error) {
	activity := models.Activity{}
	query := `UPDATE activities SET name=$1, description=$2 WHERE id=$3 RETURNING id, name, description, created_at;`
	err := db.Conn.QueryRow(query, activityData.Name, activityData.Description, activityID).Scan(&activity.ID, &activity.Name, &activity.Description, &activity.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return activity, ErrNoMatch
		}
		return activity, err
	}
	return activity, nil
}

// TODO: Archival, Restoration, Marking for Deletion
