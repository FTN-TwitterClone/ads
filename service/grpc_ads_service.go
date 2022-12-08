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

	e := model.TweetLikedEvent{
		Username: likeEvent.Username,
		TweetId:  tweetId,
	}

	err = s.eventsRepository.SaveTweetLikedEvent(serviceCtx, &e)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	//TODO: update reports

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

	e := model.TweetUnlikedEvent{
		Username: unlikeEvent.Username,
		TweetId:  tweetId,
	}

	err = s.eventsRepository.SaveTweetUnlikedEvent(serviceCtx, &e)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	//TODO: update reports

	return new(empty.Empty), nil
}
