package service

import (
	"ads/app_errors"
	"ads/model"
	"ads/repository"
	"context"
	"github.com/gocql/gocql"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type AdsService struct {
	adsRepository repository.AdsRepository
	tracer        trace.Tracer
}

func NewAdsService(adsRepository repository.AdsRepository, tracer trace.Tracer) *AdsService {
	return &AdsService{
		adsRepository,
		tracer,
	}
}

func (s *AdsService) AddProfileVisitedEvent(ctx context.Context, tweetId string, username string) *app_errors.AppError {
	serviceCtx, span := s.tracer.Start(ctx, "AdsService.AddProfileVisitedEvent")
	defer span.End()

	uuid, err := gocql.ParseUUID(tweetId)
	if err != nil {

		return &app_errors.AppError{422, "Invalid UUID"}
	}

	e := model.ProfileVisitedEvent{
		Username: username,
		TweetId:  uuid,
		Time:     time.Now(),
	}

	err = s.adsRepository.SaveProfileVisitedEvent(serviceCtx, &e)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return &app_errors.AppError{500, ""}
	}

	return nil
}

func (s *AdsService) GenerateReport(ctx context.Context, tweetId gocql.UUID, from time.Time, to time.Time) (*model.Report, *app_errors.AppError) {
	serviceCtx, span := s.tracer.Start(ctx, "AdsService.AddProfileVisitedEvent")
	defer span.End()

	visitsCount, err := s.adsRepository.GetProfileVisitsCount(serviceCtx, tweetId, from, to)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, &app_errors.AppError{500, ""}
	}

	r := model.Report{
		ProfileVisits: visitsCount,
	}

	return &r, nil
}
