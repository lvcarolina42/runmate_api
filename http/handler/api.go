package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"runmate_api/http/model"
	"runmate_api/internal/service"
)

type api struct {
	activityService *service.Activity
	userService     *service.User
}

func NewAPI(activityService *service.Activity, userService *service.User) *api {
	return &api{activityService: activityService, userService: userService}
}

func (a *api) Routes(r *chi.Mux) {
	r.Route("/activities", func(r chi.Router) {
		r.Get("/", a.getActivities)
		r.Post("/", a.createActivity)
		r.Delete("/{id}", a.deleteActivity)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", a.createUser)
		r.Get("/", a.getUsers)
		r.Get("/{id}", a.getUserByID)
		r.Get("/{username: [a-zA-Z0-9_]+}", a.getUserByUsername)
		r.Put("/{id}", a.updateUser)
		r.Delete("/{id}", a.deleteUser)

		r.Get("/{id}/activities", a.getUserActivities)
	})

	r.Post("/login", a.login)
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

func (a *api) createUser(w http.ResponseWriter, r *http.Request) {
	var input model.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := input.ToEntity()

	err = a.userService.Create(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *api) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := a.userService.ListAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := make([]*model.User, 0, len(users))
	for _, user := range users {
		result = append(result, model.NewUserFromEntity(user))
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) getUserByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := a.userService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := model.NewUserFromEntity(user)
	if result == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) getUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	user, err := a.userService.GetByUsername(r.Context(), username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := model.NewUserFromEntity(user)
	if result == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) updateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var input model.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := input.ToEntity()
	user.ID, err = uuid.Parse(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.userService.Update(r.Context(), user)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) deleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := a.userService.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) login(w http.ResponseWriter, r *http.Request) {
	var input model.LoginInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authenticated, err := a.userService.Authenticate(r.Context(), input.Username, input.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !authenticated {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
