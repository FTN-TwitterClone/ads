package service

import (
	"context"
	"fmt"
	"github.com/FTN-TwitterClone/ads/app_errors"
	"github.com/FTN-TwitterClone/ads/model"
	"github.com/FTN-TwitterClone/ads/repository"
	"github.com/gocql/gocql"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type AdsService struct {
	eventsRepository  repository.EventsRepository
	reportsRepository repository.ReportsRepository
	tracer            trace.Tracer
}

func NewAdsService(adsRepository repository.EventsRepository, reportsRepository repository.ReportsRepository, tracer trace.Tracer) *AdsService {
	return &AdsService{
		adsRepository,
		reportsRepository,
		tracer,
	}
}

func (s *AdsService) AddProfileVisitedEvent(ctx context.Context, tweetId string) *app_errors.AppError {
	serviceCtx, span := s.tracer.Start(ctx, "AdsService.AddProfileVisitedEvent")
	defer span.End()

	uuid, err := gocql.ParseUUID(tweetId)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return &app_errors.AppError{422, "Invalid UUID"}
	}

	authUser := ctx.Value("authUser").(model.AuthUser)

	now := time.Now()

	e := model.ProfileVisitedEvent{
		Username: authUser.Username,
		TweetId:  uuid,
		Time:     now,
	}

	err = s.eventsRepository.SaveProfileVisitedEvent(serviceCtx, &e)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return &app_errors.AppError{500, ""}
	}

	err = s.reportsRepository.UpsertMonthlyReportProfileVisitsCount(serviceCtx, tweetId, int64(now.Year()), int64(now.Month()))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return &app_errors.AppError{500, ""}
	}

	err = s.reportsRepository.UpsertDailyReportProfileVisitsCount(serviceCtx, tweetId, int64(now.Year()), int64(now.Month()), int64(now.Day()))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return &app_errors.AppError{500, ""}
	}

	return nil
}

func (s *AdsService) AddTweetViewedEvent(ctx context.Context, tweetId string, viewTime model.TweetViewTime) *app_errors.AppError {
	serviceCtx, span := s.tracer.Start(ctx, "AdsService.AddTweetViewedEvent")
	defer span.End()

	uuid, err := gocql.ParseUUID(tweetId)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return &app_errors.AppError{422, "Invalid UUID"}
	}

	authUser := ctx.Value("authUser").(model.AuthUser)

	now := time.Now()
	loc := now.Location()

	e := model.TweetViewedEvent{
		Username: authUser.Username,
		TweetId:  uuid,
		ViewTime: viewTime.ViewTime,
		Time:     now,
	}

	err = s.eventsRepository.SaveTweetViewedEvent(serviceCtx, &e)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return &app_errors.AppError{500, ""}
	}

	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)

	monthAvg, err := s.eventsRepository.GetAverageTweetViewTime(serviceCtx, uuid, monthStart, now)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return &app_errors.AppError{500, ""}
	}

	err = s.reportsRepository.UpsertMonthlyReportAverageProfileViewTime(serviceCtx, tweetId, int64(now.Year()), int64(now.Month()), int64(monthAvg))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return &app_errors.AppError{500, ""}
	}

	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	dayAvg, err := s.eventsRepository.GetAverageTweetViewTime(serviceCtx, uuid, dayStart, now)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return &app_errors.AppError{500, ""}
	}

	err = s.reportsRepository.UpsertDailyReportAverageProfileViewTime(serviceCtx, tweetId, int64(now.Year()), int64(now.Month()), int64(now.Day()), int64(dayAvg))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return &app_errors.AppError{500, ""}
	}

	return nil
}

func (s *AdsService) GetMonthlyReport(ctx context.Context, tweetId string, year int64, month int64) (*model.Report, *app_errors.AppError) {
	serviceCtx, span := s.tracer.Start(ctx, "AdsService.GetMonthlyReport")
	defer span.End()

	authUser := ctx.Value("authUser").(model.AuthUser)

	adInfo, err := s.eventsRepository.GetAdInfo(serviceCtx, tweetId)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, &app_errors.AppError{500, ""}
	}

	if adInfo.PostedBy != authUser.Username {
		span.SetStatus(codes.Error, fmt.Sprintf("User %s doesn't have access!", authUser.Username))
		return nil, &app_errors.AppError{403, ""}
	}

	r, err := s.reportsRepository.GetMonthlyReport(serviceCtx, tweetId, year, month)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, &app_errors.AppError{500, ""}
	}

	if r == nil {
		return &model.Report{TweetId: tweetId}, nil
	}

	return r, nil
}

func (s *AdsService) GetDailyReport(ctx context.Context, tweetId string, year int64, month int64, day int64) (*model.Report, *app_errors.AppError) {
	serviceCtx, span := s.tracer.Start(ctx, "AdsService.GetDailyReport")
	defer span.End()

	authUser := ctx.Value("authUser").(model.AuthUser)

	adInfo, err := s.eventsRepository.GetAdInfo(serviceCtx, tweetId)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, &app_errors.AppError{500, ""}
	}

	if adInfo.PostedBy != authUser.Username {
		span.SetStatus(codes.Error, fmt.Sprintf("User %s doesn't have access!", authUser.Username))
		return nil, &app_errors.AppError{403, ""}
	}

	r, err := s.reportsRepository.GetDailyReport(serviceCtx, tweetId, year, month, day)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, &app_errors.AppError{500, ""}
	}

	if r == nil {
		return &model.Report{TweetId: tweetId}, nil
	}

	return r, nil
}
