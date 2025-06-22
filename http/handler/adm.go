package handler

import (
	"encoding/json"
	"net/http"

	"runmate_api/http/model"
	"runmate_api/internal/firebase"
	"runmate_api/internal/service"

	"github.com/go-chi/chi/v5"
)

type adm struct {
	activityService  *service.Activity
	challengeService *service.Challenge
	eventService     *service.Event
	userService      *service.User

	firebaseClient *firebase.Client
}

func NewADM(
	activityService *service.Activity,
	challengeService *service.Challenge,
	eventService *service.Event,
	userService *service.User,
	firebaseClient *firebase.Client,
) *adm {
	return &adm{
		activityService:  activityService,
		challengeService: challengeService,
		eventService:     eventService,
		userService:      userService,

		firebaseClient: firebaseClient,
	}
}

func (a *adm) Routes(r *chi.Mux) {
	r.Route("/adm", func(r chi.Router) {
		r.Post("/notify", a.notify)
	})
}

func (a *adm) notify(w http.ResponseWriter, r *http.Request) {
	var input model.NotifyInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.firebaseClient.SendNotification(r.Context(), &firebase.Notification{Title: input.Title, Body: input.Body}, input.Tokens)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
