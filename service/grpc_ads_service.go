package service

import (
	"context"
	"github.com/FTN-TwitterClone/ads/repository"
	"github.com/FTN-TwitterClone/grpc-stubs/proto/ads"
	"github.com/golang/protobuf/ptypes/empty"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *gRPCAdsService) SaveAdInfo(context.Context, *ads.AdInfo) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveAdInfo not implemented")
}

func (s *gRPCAdsService) SaveLikeEvent(context.Context, *ads.LikeEvent) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveLikeEvent not implemented")
}

func (s *gRPCAdsService) SaveUnlikeEvent(context.Context, *ads.UnlikeEvent) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveUnlikeEvent not implemented")
}
