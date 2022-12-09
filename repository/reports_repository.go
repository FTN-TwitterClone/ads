package repository

import (
	"context"
	"github.com/FTN-TwitterClone/ads/model"
)

type ReportsRepository interface {
	GetMonthlyReport(ctx context.Context, tweetId string, year int64, month int64) (*model.Report, error)
	GetDailyReport(ctx context.Context, tweetId string, year int64, month int64, day int64) (*model.Report, error)
	UpsertMonthlyReportLikesCount(ctx context.Context, tweetId string, year int64, month int64) error
	UpsertMonthlyReportUnlikesCount(ctx context.Context, tweetId string, year int64, month int64) error
	UpsertMonthlyReportProfileVisitsCount(ctx context.Context, tweetId string, year int64, month int64) error
	UpsertMonthlyReportAverageProfileViewTime(ctx context.Context, tweetId string, year int64, month int64, averageViewTime int64) error
	UpsertDailyReportLikesCount(ctx context.Context, tweetId string, year int64, month int64, day int64) error
	UpsertDailyReportUnlikesCount(ctx context.Context, tweetId string, year int64, month int64, day int64) error
	UpsertDailyReportProfileVisitsCount(ctx context.Context, tweetId string, year int64, month int64, day int64) error
	UpsertDailyReportAverageProfileViewTime(ctx context.Context, tweetId string, year int64, month int64, day int64, averageViewTime int64) error
}
