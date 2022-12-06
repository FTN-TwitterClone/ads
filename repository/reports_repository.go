package repository

import (
	"context"
)

type ReportsRepository interface {
	UpsertMonthlyReportTweetLiked(ctx context.Context, tweetId string) error
}
