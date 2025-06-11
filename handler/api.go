package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"runmate_api/entity"
	"runmate_api/service"
)

type API struct {
	activityService *service.Activity
}

func NewAPI(activityService *service.Activity) *API {
	return &API{activityService: activityService}
}

func (a *API) CreateActivity(w http.ResponseWriter, r *http.Request) {
	var activity entity.Activity
	err := json.NewDecoder(r.Body).Decode(&activity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.activityService.Create(r.Context(), &activity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *API) GetActivities(w http.ResponseWriter, r *http.Request) {
	activities, err := a.activityService.ListAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(activities)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) GetActivitiesByUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	activities, err := a.activityService.ListByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(activities)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) DeleteActivity(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := a.activityService.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
