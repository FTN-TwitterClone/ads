package controller

import (
	"ads/service"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

type AdsController struct {
	tweetService *service.AdsService
	tracer       trace.Tracer
}

func NewAdsController(tweetService *service.AdsService, tracer trace.Tracer) *AdsController {
	return &AdsController{
		tweetService,
		tracer,
	}
}

func (c *AdsController) DeleteLike(w http.ResponseWriter, req *http.Request) {

}
