package service

import (
	"context"
	"sort"

	"runmate_api/internal/entity"
	"runmate_api/internal/repository"
)

type Activity struct {
	activityRepo  *repository.Activity
	challengeRepo *repository.Challenge
	userRepo      *repository.User
}

func NewActivity(activityRepo *repository.Activity, challengeRepo *repository.Challenge, userRepo *repository.User) *Activity {
	return &Activity{activityRepo: activityRepo, challengeRepo: challengeRepo, userRepo: userRepo}
}

func (a *Activity) Create(ctx context.Context, activity *entity.Activity) error {
	user, err := a.userRepo.GetByID(ctx, activity.UserID.String())
	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	for i, coordinate := range activity.Coordinates {
		coordinate.Order = int(i)
	}

	err = a.activityRepo.Create(ctx, activity)
	if err != nil {
		return err
	}

	user.XP += activity.Distance
	err = a.userRepo.Update(ctx, user)
	if err != nil {
		return err
	}

	userChallenges, err := a.challengeRepo.GetAllActiveByUser(ctx, user)
	if err != nil {
		return err
	}

	for _, userChallenge := range userChallenges {
		if activity.Date.Before(userChallenge.StartDate) || (userChallenge.EndDate != nil && activity.Date.After(*userChallenge.EndDate)) {
			continue
		}

		err = a.challengeRepo.AddEvent(ctx, userChallenge, &entity.ChallengeEvent{
			ChallengeID: userChallenge.ID,
			UserID:      user.ID,
			Distance:    activity.Distance,
			Date:        activity.Date,
		})
		if err != nil {
			return err
		}

		if userChallenge.Type == entity.ChallengeTypeDistance {
			userChallengeEvents, err := a.challengeRepo.GetAllEventsByUser(ctx, userChallenge, user)
			if err != nil {
				return err
			}

			var total int
			for _, userChallengeEvent := range userChallengeEvents {
				total += userChallengeEvent.Distance
			}

			if total >= *userChallenge.TotalDistance {
				userChallenge.EndDate = &activity.Date
				err = a.challengeRepo.Update(ctx, userChallenge)
				if err != nil {
					return err
				}
			}
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
