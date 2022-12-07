package repository

import (
	"context"
	"github.com/FTN-TwitterClone/ads/model"
	"github.com/gocql/gocql"
	"time"
)

type EventsRepository interface {
	SaveAdInfo(ctx context.Context, adInfo *model.AdInfo) error
	SaveTweetLikedEvent(ctx context.Context, tweetLikedEvent *model.TweetLikedEvent) error
	SaveTweetUnlikedEvent(ctx context.Context, tweetUnlikedEvent *model.TweetUnlikedEvent) error
	SaveTweetViewedEvent(ctx context.Context, tweetViewedEvent *model.TweetViewedEvent) error
	SaveProfileVisitedEvent(ctx context.Context, profileVisitedEvent *model.ProfileVisitedEvent) error
	GetAverageTweetViewTime(ctx context.Context, tweetId gocql.UUID, from time.Time, to time.Time) (int, error)
}
