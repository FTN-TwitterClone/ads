package service

import (
	"context"
	"github.com/FTN-TwitterClone/ads/model"
	"github.com/FTN-TwitterClone/ads/repository"
	"github.com/FTN-TwitterClone/grpc-stubs/proto/ads"
	"github.com/gocql/gocql"
	"github.com/golang/protobuf/ptypes/empty"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type gRPCAdsService struct {
	ads.UnimplementedAdsServiceServer
	tracer            trace.Tracer
	eventsRepository  repository.EventsRepository
	reportsRepository repository.ReportsRepository
}

func NewgRPCAdsService(tracer trace.Tracer, eventsRepository repository.EventsRepository, reportsRepository repository.ReportsRepository) *gRPCAdsService {
	return &gRPCAdsService{
		tracer:            tracer,
		eventsRepository:  eventsRepository,
		reportsRepository: reportsRepository,
	}
}

func (s *gRPCAdsService) SaveAdInfo(ctx context.Context, adInfo *ads.AdInfo) (*empty.Empty, error) {
	serviceCtx, span := s.tracer.Start(ctx, "gRPCAdsService.SaveAdInfo")
	defer span.End()

	tweetId, err := gocql.ParseUUID(adInfo.TweetId)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	a := model.AdInfo{
		TweetId:  tweetId,
		PostedBy: adInfo.PostedBy,
		Town:     adInfo.Town,
		MinAge:   adInfo.MinAge,
		MaxAge:   adInfo.MaxAge,
		Gender:   adInfo.Gender,
	}

	err = s.eventsRepository.SaveAdInfo(serviceCtx, &a)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return new(empty.Empty), nil
}

func (s *gRPCAdsService) SaveLikeEvent(ctx context.Context, likeEvent *ads.LikeEvent) (*empty.Empty, error) {
	serviceCtx, span := s.tracer.Start(ctx, "gRPCAdsService.SaveLikeEvent")
	defer span.End()

	tweetId, err := gocql.ParseUUID(likeEvent.TweetId)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	now := time.Now()

	e := model.TweetLikedEvent{
		Username: likeEvent.Username,
		TweetId:  tweetId,
		Time:     now,
	}

	err = s.eventsRepository.SaveTweetLikedEvent(serviceCtx, &e)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	err = s.reportsRepository.UpsertMonthlyReportLikesCount(serviceCtx, likeEvent.TweetId, int64(now.Year()), int64(now.Month()))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	err = s.reportsRepository.UpsertDailyReportLikesCount(serviceCtx, likeEvent.TweetId, int64(now.Year()), int64(now.Month()), int64(now.Day()))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return new(empty.Empty), nil
}

func (s *gRPCAdsService) SaveUnlikeEvent(ctx context.Context, unlikeEvent *ads.UnlikeEvent) (*empty.Empty, error) {
	serviceCtx, span := s.tracer.Start(ctx, "gRPCAdsService.SaveLikeEvent")
	defer span.End()

	tweetId, err := gocql.ParseUUID(unlikeEvent.TweetId)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	now := time.Now()

	e := model.TweetUnlikedEvent{
		Username: unlikeEvent.Username,
		TweetId:  tweetId,
		Time:     now,
	}

	err = s.eventsRepository.SaveTweetUnlikedEvent(serviceCtx, &e)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	err = s.reportsRepository.UpsertMonthlyReportUnlikesCount(serviceCtx, unlikeEvent.TweetId, int64(now.Year()), int64(now.Month()))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	err = s.reportsRepository.UpsertDailyReportUnlikesCount(serviceCtx, unlikeEvent.TweetId, int64(now.Year()), int64(now.Month()), int64(now.Day()))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return new(empty.Empty), nil
}
