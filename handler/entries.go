package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/turbo-d/overuse/db"
	"github.com/turbo-d/overuse/models"
)

var entryIDKey = "entryID"

func entries(router chi.Router) {
	router.Get("/", getAllEntries)
	router.Post("/", createEntry)
	router.Route("/{entryID}", func(router chi.Router) {
		router.Use(EntryContext)
		router.Get("/", getEntry)
		router.Put("/", updateEntry)
		router.Delete("/", deleteEntry)
	})
}

func EntryContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entryID := chi.URLParam(r, "entryID")
		if entryID == "" {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("entry ID is required")))
			return
		}
		id, err := strconv.Atoi(entryID)
		if err != nil {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("invalid entry ID")))
		}
		ctx := context.WithValue(r.Context(), entryIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func createEntry(w http.ResponseWriter, r *http.Request) {
	entry := &models.Entry{}
	if err := render.Bind(r, entry); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	if err := dbInstance.AddEntry(entry); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, entry); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func getAllEntries(w http.ResponseWriter, r *http.Request) {
	entries, err := dbInstance.GetAllEntries()
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, entries); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
}

func getEntry(w http.ResponseWriter, r *http.Request) {
	entryID := r.Context().Value(entryIDKey).(int)
	entry, err := dbInstance.GetEntryById(entryID)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &entry); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func deleteEntry(w http.ResponseWriter, r *http.Request) {
	entryID := r.Context().Value(entryIDKey).(int)
	err := dbInstance.DeleteEntry(entryID)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
}

func updateEntry(w http.ResponseWriter, r *http.Request) {
	entryID := r.Context().Value(entryIDKey).(int)
	entryData := models.Entry{}
	if err := render.Bind(r, &entryData); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	entry, err := dbInstance.UpdateEntry(entryID, entryData)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &entry); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}
