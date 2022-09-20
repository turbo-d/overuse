package models

import (
	"net/http"
	"time"
)

type Entry struct {
	ID         int       `json:"id"`
	ActivityID int       `json:"activity_id"`
	OccuredOn  time.Time `json:"occured_on"`
	Volume     int       `json:"volume"`
	CreatedAt  time.Time `json:"created_at"`
}

type EntryList struct {
	Entries []Entry `json:"entries"`
}

func (e *Entry) Bind(r *http.Request) error {
	return nil
}

func (*EntryList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (*Entry) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
