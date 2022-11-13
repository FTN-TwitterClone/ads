package service

import (
	"ads/app_errors"
	"ads/model"
	"ads/repository"
	"context"
	"go.opentelemetry.io/otel/trace"
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

func (s *AdsService) CreateTweet(ctx context.Context, tweet model.Tweet) (*model.Tweet, *app_errors.AppError) {
	_, span := s.tracer.Start(ctx, "AdsService.CreateTweet")
	defer span.End()

	return nil, nil
}
