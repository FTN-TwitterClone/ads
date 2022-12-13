package controller

import (
	"fmt"
	"github.com/FTN-TwitterClone/ads/controller/json"
	"github.com/FTN-TwitterClone/ads/model"
	"github.com/FTN-TwitterClone/ads/service"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strconv"
)

type AdsController struct {
	adsService *service.AdsService
	tracer     trace.Tracer
}

func NewAdsController(tweetService *service.AdsService, tracer trace.Tracer) *AdsController {
	return &AdsController{
		tweetService,
		tracer,
	}
}

func (c *AdsController) GetAdInfo(w http.ResponseWriter, req *http.Request) {
	ctx, span := c.tracer.Start(req.Context(), "AdsController.GetAdInfo")
	defer span.End()

	authUser := ctx.Value("authUser").(model.AuthUser)

	if authUser.Role != "ROLE_BUSINESS" {
		span.SetStatus(codes.Error, fmt.Sprintf("%s not allowed!", authUser.Role))
		http.Error(w, "", 403)
		return
	}

	tweetId := mux.Vars(req)["tweetId"]

	adInfo, appErr := c.adsService.GetAdInfo(ctx, tweetId)
	if appErr != nil {
		span.SetStatus(codes.Error, appErr.Error())
		http.Error(w, appErr.Message, appErr.Code)
		return
	}

	json.EncodeJson(w, &adInfo)
}

func (c *AdsController) AddProfileVisitedEvent(w http.ResponseWriter, req *http.Request) {
	ctx, span := c.tracer.Start(req.Context(), "AdsController.AddProfileVisitedEvent")
	defer span.End()

	tweetId := mux.Vars(req)["tweetId"]

	appErr := c.adsService.AddProfileVisitedEvent(ctx, tweetId)
	if appErr != nil {
		span.SetStatus(codes.Error, appErr.Error())
		http.Error(w, appErr.Message, appErr.Code)
		return
	}
}

func (c *AdsController) AddTweetViewedEvent(w http.ResponseWriter, req *http.Request) {
	ctx, span := c.tracer.Start(req.Context(), "AdsController.AddTweetViewedEvent")
	defer span.End()

	tweetId := mux.Vars(req)["tweetId"]

	viewTime, err := json.DecodeJson[model.TweetViewTime](req.Body)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	appErr := c.adsService.AddTweetViewedEvent(ctx, tweetId, viewTime)
	if appErr != nil {
		span.SetStatus(codes.Error, appErr.Error())
		http.Error(w, appErr.Message, appErr.Code)
		return
	}
}

func (c *AdsController) GetMonthlyReport(w http.ResponseWriter, req *http.Request) {
	ctx, span := c.tracer.Start(req.Context(), "AdsController.GetMonthlyReport")
	defer span.End()

	authUser := ctx.Value("authUser").(model.AuthUser)

	if authUser.Role != "ROLE_BUSINESS" {
		span.SetStatus(codes.Error, fmt.Sprintf("%s not allowed!", authUser.Role))
		http.Error(w, "", 403)
		return
	}

	tweetIdString := mux.Vars(req)["tweetId"]

	year, err := strconv.ParseInt(mux.Vars(req)["year"], 10, 64)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "", 500)
		return
	}

	month, err := strconv.ParseInt(mux.Vars(req)["month"], 10, 64)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "", 500)
		return
	}

	tweetId, err := gocql.ParseUUID(tweetIdString)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "", 500)
		return
	}

	report, appErr := c.adsService.GetMonthlyReport(ctx, tweetId.String(), year, month)
	if appErr != nil {
		span.SetStatus(codes.Error, appErr.Error())
		http.Error(w, appErr.Message, appErr.Code)
		return
	}

	json.EncodeJson(w, &report)
}

func (c *AdsController) GetDailyReport(w http.ResponseWriter, req *http.Request) {
	ctx, span := c.tracer.Start(req.Context(), "AdsController.GetDailyReport")
	defer span.End()

	authUser := ctx.Value("authUser").(model.AuthUser)

	if authUser.Role != "ROLE_BUSINESS" {
		span.SetStatus(codes.Error, fmt.Sprintf("%s not allowed!", authUser.Role))
		http.Error(w, "", 403)
		return
	}

	tweetIdString := mux.Vars(req)["tweetId"]

	year, err := strconv.ParseInt(mux.Vars(req)["year"], 10, 64)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "", 500)
		return
	}

	month, err := strconv.ParseInt(mux.Vars(req)["month"], 10, 64)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "", 500)
		return
	}

	day, err := strconv.ParseInt(mux.Vars(req)["day"], 10, 64)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "", 500)
		return
	}

	tweetId, err := gocql.ParseUUID(tweetIdString)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "", 500)
		return
	}

	report, appErr := c.adsService.GetDailyReport(ctx, tweetId.String(), year, month, day)
	if appErr != nil {
		span.SetStatus(codes.Error, appErr.Error())
		http.Error(w, appErr.Message, appErr.Code)
		return
	}

	json.EncodeJson(w, &report)
}
