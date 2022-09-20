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

var activityIDKey = "activityID"

func activities(router chi.Router) {
	router.Get("/", getAllActivities)
	router.Post("/", createActivity)
	router.Route("/{activityID}", func(router chi.Router) {
		router.Use(ActivityContext)
		router.Get("/", getActivity)
		router.Put("/", updateActivity)
		router.Delete("/", deleteActivity)
	})
}

func ActivityContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		activityID := chi.URLParam(r, "activityID")
		if activityID == "" {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("activity ID is required")))
			return
		}
		id, err := strconv.Atoi(activityID)
		if err != nil {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("invalid activity ID")))
		}
		ctx := context.WithValue(r.Context(), activityIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func createActivity(w http.ResponseWriter, r *http.Request) {
	activity := &models.Activity{}
	if err := render.Bind(r, activity); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	if err := dbInstance.AddActivity(activity); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, activity); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func getAllActivities(w http.ResponseWriter, r *http.Request) {
	activities, err := dbInstance.GetAllActivities()
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, activities); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
}

func getActivity(w http.ResponseWriter, r *http.Request) {
	activityID := r.Context().Value(activityIDKey).(int)
	activity, err := dbInstance.GetActivityById(activityID)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &activity); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func deleteActivity(w http.ResponseWriter, r *http.Request) {
	activityID := r.Context().Value(activityIDKey).(int)
	err := dbInstance.DeleteActivity(activityID)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
}

func updateActivity(w http.ResponseWriter, r *http.Request) {
	activityID := r.Context().Value(activityIDKey).(int)
	activityData := models.Activity{}
	if err := render.Bind(r, &activityData); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	activity, err := dbInstance.UpdateActivity(activityID, activityData)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &activity); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}
