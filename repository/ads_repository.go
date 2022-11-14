package repository

import (
	"ads/model"
	"context"
	"time"
)

type AdsRepository interface {
	SaveTweetLikedEvent(ctx context.Context, tweetLikedEvent *model.TweetLikedEvent) error
	SaveProfileVisitedEvent(ctx context.Context, profileVisitedEvent *model.ProfileVisitedEvent) error
	GetProfileVisitsCount(ctx context.Context, from time.Time, to time.Time) (int, error)
}
