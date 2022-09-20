package models

import (
	"fmt"
	"net/http"
	"time"
)

type Activity struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Units       string    `json:"units"`
	IsDeleted   bool      `json:"is_deleted"`
	IsArchived  bool      `json:"is_archived"`
	CreatedAt   time.Time `json:"created_at"`
}

type ActivityList struct {
	Activities []Activity `json:"activities"`
}

func (a *Activity) Bind(r *http.Request) error {
	if a.Name == "" {
		return fmt.Errorf("name is a required field")
	}
	if a.Units == "" {
		return fmt.Errorf("units is a required field")
	}
	return nil
}

func (*ActivityList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (*Activity) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
