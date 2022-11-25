package repository

import (
	"context"
	"github.com/FTN-TwitterClone/ads/model"
	"github.com/gocql/gocql"
	"time"
)

type AdsRepository interface {
	SaveTweetLikedEvent(ctx context.Context, tweetLikedEvent *model.TweetLikedEvent) error
	SaveTweetUnlikedEvent(ctx context.Context, tweetUnlikedEvent *model.TweetUnlikedEvent) error
	SaveTweetViewedEvent(ctx context.Context, tweetViewedEvent *model.TweetViewedEvent) error
	SaveProfileVisitedEvent(ctx context.Context, profileVisitedEvent *model.ProfileVisitedEvent) error
	GetTweetLikesCount(ctx context.Context, tweetId gocql.UUID, from time.Time, to time.Time) (int, error)
	GetTweetUnlikesCount(ctx context.Context, tweetId gocql.UUID, from time.Time, to time.Time) (int, error)
	GetAverageTweetViewTimeCount(ctx context.Context, tweetId gocql.UUID, from time.Time, to time.Time) (int, error)
	GetProfileVisitsCount(ctx context.Context, tweetId gocql.UUID, from time.Time, to time.Time) (int, error)
}
