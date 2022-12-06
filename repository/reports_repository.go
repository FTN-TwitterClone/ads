package repository

import (
	"context"
)

type ReportsRepository interface {
	UpsertMonthlyReportLikesCount(ctx context.Context, tweetId string, year int64, month int64) error
	UpsertMonthlyReportAverageProfileViewTime(ctx context.Context, tweetId string, year int64, month int64, averageViewTime int64) error
}
