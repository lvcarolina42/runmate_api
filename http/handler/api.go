package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"runmate_api/http/model"
	"runmate_api/internal/service"
)

type api struct {
	activityService *service.Activity
}

func NewAPI(activityService *service.Activity) *api {
	return &api{activityService: activityService}
}

func (a *api) Routes(r *chi.Mux) {
	r.Route("/v2", func(r chi.Router) {
		r.Route("/activities", func(r chi.Router) {
			r.Get("/", a.getActivities)
			r.Post("/", a.createActivity)
			r.Delete("/{id}", a.deleteActivity)
		})

		r.Route("/users", func(r chi.Router) {
			//r.Post("/", a.createUser)
			//r.Get("/", a.getUsers)
			//r.Get("/{id}", a.getUser)
			//r.Put("/{id}", a.updateUser)
			//r.Delete("/{id}", a.deleteUser)

			r.Get("/{id}/activities", a.getUserActivities)
		})
	})
}

func (a *api) createActivity(w http.ResponseWriter, r *http.Request) {
	var input model.Activity
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	activity, err := input.ToEntity()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.activityService.Create(r.Context(), activity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *api) getActivities(w http.ResponseWriter, r *http.Request) {
	activities, err := a.activityService.ListAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := make([]*model.Activity, 0, len(activities))
	for _, activity := range activities {
		result = append(result, model.NewActivityFromEntity(activity))
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) deleteActivity(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := a.activityService.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) getUserActivities(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	activities, err := a.activityService.ListByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := make([]*model.Activity, 0, len(activities))
	for _, activity := range activities {
		result = append(result, model.NewActivityFromEntity(activity))
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
