package service

import (
	"context"
	"maps"
	"slices"
	"sort"

	"runmate_api/internal/entity"
	"runmate_api/internal/firebase"
	"runmate_api/internal/repository"
)

type Activity struct {
	activityRepo  *repository.Activity
	challengeRepo *repository.Challenge
	userRepo      *repository.User

	firebaseClient *firebase.Client
}

func NewActivity(activityRepo *repository.Activity, challengeRepo *repository.Challenge, userRepo *repository.User, firebaseClient *firebase.Client) *Activity {
	return &Activity{
		activityRepo:  activityRepo,
		challengeRepo: challengeRepo,
		userRepo:      userRepo,

		firebaseClient: firebaseClient,
	}
}

func (a *Activity) Create(ctx context.Context, activity *entity.Activity) error {
	owner, err := a.userRepo.GetByID(ctx, activity.UserID.String())
	if err != nil {
		return err
	}

	if owner == nil {
		return ErrUserNotFound
	}

	for i, coordinate := range activity.Coordinates {
		coordinate.Order = int(i)
	}

	err = a.activityRepo.Create(ctx, activity)
	if err != nil {
		return err
	}

	owner.XP += activity.Distance
	err = a.userRepo.Update(ctx, owner)
	if err != nil {
		return err
	}

	ownerChallenges, err := a.challengeRepo.GetAllActiveByUser(ctx, owner)
	if err != nil {
		return err
	}

	for _, ownerChallenge := range ownerChallenges {
		if activity.Date.Before(ownerChallenge.StartDate) || (ownerChallenge.EndDate != nil && activity.Date.After(*ownerChallenge.EndDate)) {
			continue
		}

		err = a.challengeRepo.AddEvent(ctx, ownerChallenge, &entity.ChallengeEvent{
			ChallengeID: ownerChallenge.ID,
			UserID:      owner.ID,
			Distance:    activity.Distance,
			Date:        activity.Date,
		})
		if err != nil {
			return err
		}

		tokens := make(map[string]any, len(ownerChallenge.Users)-1)
		for _, user := range ownerChallenge.Users {
			if user.ID == owner.ID || owner.FCMToken == "" {
				continue
			}

			tokens[owner.FCMToken] = struct{}{}
		}

		notificationFunc := newChallengeActivityNotification
		if ownerChallenge.Type == entity.ChallengeTypeDistance {
			userChallengeEvents, err := a.challengeRepo.GetAllEventsByUser(ctx, ownerChallenge, owner)
			if err != nil {
				return err
			}

			var total int
			for _, userChallengeEvent := range userChallengeEvents {
				total += userChallengeEvent.Distance
			}

			if total >= *ownerChallenge.TotalDistance {
				ownerChallenge.EndDate = &activity.Date
				err = a.challengeRepo.Update(ctx, ownerChallenge)
				if err != nil {
					return err
				}

				notificationFunc = endChallengeNotification
			}
		}

		notification := notificationFunc(owner.Name, ownerChallenge.Title)
		err = a.firebaseClient.SendNotification(ctx, notification, slices.Collect(maps.Keys(tokens)))
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *Activity) ListAll(ctx context.Context) ([]*entity.Activity, error) {
	return a.activityRepo.GetAll(ctx)
}

func (a *Activity) ListByUser(ctx context.Context, userID string) ([]*entity.Activity, error) {
	return a.activityRepo.GetByUserID(ctx, userID)
}

func (a *Activity) ListAllFromUserFriends(ctx context.Context, userID string) ([]*entity.Activity, error) {
	user, err := a.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	friends, err := a.userRepo.ListFriends(ctx, user)
	if err != nil {
		return nil, err
	}

	var activities []*entity.Activity
	for _, friend := range friends {
		activitiesFromFriend, err := a.activityRepo.GetByUserID(ctx, friend.ID.String())
		if err != nil {
			return nil, err
		}

		activities = append(activities, activitiesFromFriend...)
	}

	sort.Slice(activities, func(i, j int) bool {
		return activities[i].Date.Before(activities[j].Date)
	})

	return activities, nil
}

func (a *Activity) Delete(ctx context.Context, id string) error {
	return a.activityRepo.Delete(ctx, id)
}
