package service

import (
	"context"
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
	}

	err = s.eventsRepository.SaveProfileVisitedEvent(serviceCtx, &e)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return &app_errors.AppError{500, ""}
	}

	return nil
}

func (s *AdsService) GenerateReport(ctx context.Context, tweetId gocql.UUID, from time.Time, to time.Time) (*model.Report, *app_errors.AppError) {
	//serviceCtx, span := s.tracer.Start(ctx, "AdsService.AddProfileVisitedEvent")
	//defer span.End()
	//
	//likesCount, err := s.eventsRepository.GetTweetLikesCount(serviceCtx, tweetId, from, to)
	//if err != nil {
	//	span.SetStatus(codes.Error, err.Error())
	//	return nil, &app_errors.AppError{500, ""}
	//}
	//
	//unlikesCount, err := s.eventsRepository.GetTweetUnlikesCount(serviceCtx, tweetId, from, to)
	//if err != nil {
	//	span.SetStatus(codes.Error, err.Error())
	//	return nil, &app_errors.AppError{500, ""}
	//}
	//
	//viewTime, err := s.eventsRepository.GetAverageTweetViewTimeCount(serviceCtx, tweetId, from, to)
	//if err != nil {
	//	span.SetStatus(codes.Error, err.Error())
	//	return nil, &app_errors.AppError{500, ""}
	//}
	//
	//visitsCount, err := s.eventsRepository.GetProfileVisitsCount(serviceCtx, tweetId, from, to)
	//if err != nil {
	//	span.SetStatus(codes.Error, err.Error())
	//	return nil, &app_errors.AppError{500, ""}
	//}
	//
	//r := model.Report{
	//	TweetsLiked:          likesCount,
	//	TweetsUnliked:        unlikesCount,
	//	AverageTweetViewTime: viewTime,
	//	ProfileVisits:        visitsCount,
	//	From:                 from,
	//	To:                   to,
	//}
	//
	//return &r, nil

	return nil, nil
}
