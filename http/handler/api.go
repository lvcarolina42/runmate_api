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
	activityService  *service.Activity
	challengeService *service.Challenge
	userService      *service.User
}

func NewAPI(
	activityService *service.Activity,
	challengeService *service.Challenge,
	userService *service.User,
) *api {
	return &api{
		activityService:  activityService,
		challengeService: challengeService,
		userService:      userService,
	}
}

func (a *api) Routes(r *chi.Mux) {
	r.Route("/activities", func(r chi.Router) {
		r.Get("/", a.getActivities)
		r.Post("/", a.createActivity)
		r.Delete("/{id}", a.deleteActivity)
	})

	r.Route("/challenges", func(r chi.Router) {
		r.Post("/", a.createChallenge)
		r.Get("/{id}", a.getChallenge)
		r.Put("/join", a.joinChallenge)
	})

	r.Route("/friends", func(r chi.Router) {
		r.Post("/", a.addFriend)
		r.Delete("/", a.removeFriend)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", a.createUser)
		r.Get("/", a.getUsers)
		r.Get("/{username:[a-zA-Z0-9_]+}", a.getUserByUsername)
		r.Get("/{id:[a-zA-Z0-9\\-]{36}}", a.getUserByID)
		r.Put("/{id}", a.updateUser)
		r.Delete("/{id}", a.deleteUser)

		r.Get("/{id}/activities", a.getUserActivities)

		r.Get("/{id}/challenges", a.getUserChallenges)

		r.Get("/{id}/friends", a.listFriends)
	})

	r.Post("/login", a.login)
}

func (a *api) createActivity(w http.ResponseWriter, r *http.Request) {
	var input model.CreateActivityInput
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

func (a *api) createChallenge(w http.ResponseWriter, r *http.Request) {
	var input *model.CreateChallengeInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = input.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	challenge, err := input.ToEntity()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.challengeService.Create(r.Context(), challenge)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *api) getChallenge(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	challenge, err := a.challengeService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if challenge == nil {
		http.Error(w, "challenge not found", http.StatusNotFound)
		return
	}

	ranking, err := a.challengeService.GetRanking(r.Context(), challenge)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(model.NewChallengeFromEntity(challenge, ranking))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) joinChallenge(w http.ResponseWriter, r *http.Request) {
	var input model.JoinChallengeInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.challengeService.Join(r.Context(), input.ChallengeID, input.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) getUserChallenges(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	listChallengesFunc := a.challengeService.ListAllByUserID
	if r.URL.Query().Get("active") == "1" {
		listChallengesFunc = a.challengeService.ListAllActiveByUserID
	}

	challenges, err := listChallengesFunc(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := make([]*model.Challenge, 0, len(challenges))
	for _, challenge := range challenges {
		ranking, err := a.challengeService.GetRanking(r.Context(), challenge)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result = append(result, model.NewChallengeFromEntity(challenge, ranking))
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

func (a *api) addFriend(w http.ResponseWriter, r *http.Request) {
	var input model.FriendInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.userService.AddFriend(r.Context(), input.UserID, input.FriendID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) listFriends(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	friends, err := a.userService.ListFriends(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := make([]*model.User, 0, len(friends))
	for _, friend := range friends {
		result = append(result, model.NewUserFromEntity(friend))
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) removeFriend(w http.ResponseWriter, r *http.Request) {
	var input model.FriendInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.userService.RemoveFriend(r.Context(), input.UserID, input.FriendID)
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

	user, err := a.userService.Authenticate(r.Context(), input.Username, input.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err = json.NewEncoder(w).Encode(model.NewUserFromEntity(user))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
