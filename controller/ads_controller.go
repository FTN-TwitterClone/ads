package controller

import (
	"github.com/FTN-TwitterClone/ads/controller/json"
	"github.com/FTN-TwitterClone/ads/service"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strconv"
	"time"
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

func (c *AdsController) AddProfileVisitedEvent(w http.ResponseWriter, req *http.Request) {
	ctx, span := c.tracer.Start(req.Context(), "AdsController.AddProfileVisitedEvent")
	defer span.End()

	username := mux.Vars(req)["username"]
	tweetId := mux.Vars(req)["tweetId"]

	appErr := c.adsService.AddProfileVisitedEvent(ctx, tweetId, username)
	if appErr != nil {
		span.SetStatus(codes.Error, appErr.Error())
		http.Error(w, appErr.Message, appErr.Code)
		return
	}
}

func (c *AdsController) GetMonthlyReport(w http.ResponseWriter, req *http.Request) {
	ctx, span := c.tracer.Start(req.Context(), "AdsController.GetReport")
	defer span.End()

	tweetIdString := mux.Vars(req)["id"]

	fromInt, err := strconv.ParseInt(mux.Vars(req)["from"], 10, 64)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "", 500)
		return
	}

	toInt, err := strconv.ParseInt(mux.Vars(req)["to"], 10, 64)
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

	from := time.UnixMilli(fromInt)
	to := time.UnixMilli(toInt)

	report, appErr := c.adsService.GenerateReport(ctx, tweetId, from, to)
	if appErr != nil {
		span.SetStatus(codes.Error, appErr.Error())
		http.Error(w, appErr.Message, appErr.Code)
		return
	}

	json.EncodeJson(w, &report)
}

func (c *AdsController) GetDailyReport(w http.ResponseWriter, req *http.Request) {

}
