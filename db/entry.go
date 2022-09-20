package db

import (
	"database/sql"
	"time"

	"github.com/turbo-d/overuse/models"
)

func (db Database) GetAllEntries() (*models.EntryList, error) {
	list := &models.EntryList{}
	rows, err := db.Conn.Query("SELECT * FROM entries ORDER BY ID DESC")
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var entry models.Entry
		err := rows.Scan(&entry.ID, &entry.ActivityID, &entry.OccuredOn, &entry.Volume, &entry.CreatedAt)
		if err != nil {
			return list, err
		}
		list.Entries = append(list.Entries, entry)
	}
	return list, nil
}

func (db Database) AddEntry(entry *models.Entry) error {
	var id int
	var createdAt time.Time
	query := `INSERT INTO entries (activity_id, occured_on, volume) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := db.Conn.QueryRow(query, entry.ActivityID, entry.OccuredOn, entry.Volume).Scan(&id, &createdAt)
	if err != nil {
		return err
	}
	entry.ID = id
	entry.CreatedAt = createdAt
	return nil
}

func (db Database) GetEntryById(entryID int) (models.Entry, error) {
	entry := models.Entry{}
	query := `SELECT * FROM entries WHERE id = $1;`
	row := db.Conn.QueryRow(query, entryID)
	switch err := row.Scan(&entry.ID, &entry.ActivityID, &entry.OccuredOn, &entry.Volume, &entry.CreatedAt); err {
	case sql.ErrNoRows:
		return entry, ErrNoMatch
	default:
		return entry, err
	}
}

func (db Database) DeleteEntry(entryID int) error {
	query := `DELETE FROM entries WHERE id = $1;`
	_, err := db.Conn.Exec(query, entryID)
	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}

func (db Database) UpdateEntry(entryID int, entryData models.Entry) (models.Entry, error) {
	entry := models.Entry{}
	query := `UPDATE entries SET occured_on=$1, volume=$2 WHERE id=$3 RETURNING id, occured_on, volume, created_at;`
	err := db.Conn.QueryRow(query, entryData.OccuredOn, entryData.Volume, entryID).Scan(&entry.ID, &entry.OccuredOn, &entry.Volume, &entry.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return entry, ErrNoMatch
		}
		return entry, err
	}
	return entry, nil
}
