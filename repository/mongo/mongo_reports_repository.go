package mongo

import (
	"context"
	"fmt"
	"github.com/FTN-TwitterClone/ads/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"os"
)

const (
	MONTHLY = "monthly"
	DAILY   = "daily"
)

type MongoReportsRepository struct {
	tracer trace.Tracer
	cli    *mongo.Client
}

func NewMongoReportsRepository(tracer trace.Tracer) (*MongoReportsRepository, error) {

	db := os.Getenv("MONGO_DB")
	dbport := os.Getenv("MONGO_DBPORT")

	//mongo logic
	host := fmt.Sprintf("%s:%s", db, dbport)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(`mongodb://`+host))
	if err != nil {
		panic(err)
	}
	client.Database("reportsDB").Collection("reports")

	car := MongoReportsRepository{
		tracer,
		client,
	}

	return &car, nil
}

func (r *MongoReportsRepository) GetMonthlyReport(ctx context.Context, tweetId string, year int64, month int64) (*model.Report, error) {
	_, span := r.tracer.Start(ctx, "MongoReportsRepository.GetMonthlyReport")
	defer span.End()

	usersCollection := r.cli.Database("reportsDB").Collection("reports")

	filter := bson.M{"tweetId": tweetId, "type": MONTHLY, "year": year, "month": month}

	var report model.Report

	res := usersCollection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return nil, nil
	}

	err := res.Decode(&report)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &report, nil
}

func (r *MongoReportsRepository) GetDailyReport(ctx context.Context, tweetId string, year int64, month int64, day int64) (*model.Report, error) {
	_, span := r.tracer.Start(ctx, "MongoReportsRepository.GetMonthlyReport")
	defer span.End()

	usersCollection := r.cli.Database("reportsDB").Collection("reports")

	filter := bson.M{"tweetId": tweetId, "type": DAILY, "year": year, "month": month, "day": day}

	var report model.Report

	res := usersCollection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return nil, nil
	}

	err := res.Decode(&report)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &report, nil
}

func (r *MongoReportsRepository) UpsertMonthlyReportLikesCount(ctx context.Context, tweetId string, year int64, month int64) error {
	_, span := r.tracer.Start(ctx, "MongoReportsRepository.UpsertMonthlyReportLikesCount")
	defer span.End()

	usersCollection := r.cli.Database("reportsDB").Collection("reports")

	filter := bson.M{"tweetId": tweetId, "type": MONTHLY, "year": year, "month": month}
	update := bson.D{{"$inc", bson.D{{"likesCount", 1}}}}
	setUpsert := options.Update().SetUpsert(true)

	_, err := usersCollection.UpdateOne(ctx, filter, update, setUpsert)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func (r *MongoReportsRepository) UpsertMonthlyReportUnlikesCount(ctx context.Context, tweetId string, year int64, month int64) error {
	_, span := r.tracer.Start(ctx, "MongoReportsRepository.UpsertMonthlyReportUnlikesCount")
	defer span.End()

	usersCollection := r.cli.Database("reportsDB").Collection("reports")

	filter := bson.M{"tweetId": tweetId, "type": MONTHLY, "year": year, "month": month}
	update := bson.D{{"$inc", bson.D{{"unlikesCount", 1}}}}
	setUpsert := options.Update().SetUpsert(true)

	_, err := usersCollection.UpdateOne(ctx, filter, update, setUpsert)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func (r *MongoReportsRepository) UpsertMonthlyReportProfileVisitsCount(ctx context.Context, tweetId string, year int64, month int64) error {
	_, span := r.tracer.Start(ctx, "MongoReportsRepository.UpsertMonthlyReportProfileVisitsCount")
	defer span.End()

	usersCollection := r.cli.Database("reportsDB").Collection("reports")

	filter := bson.M{"tweetId": tweetId, "type": MONTHLY, "year": year, "month": month}
	update := bson.D{{"$inc", bson.D{{"profileVisits", 1}}}}
	setUpsert := options.Update().SetUpsert(true)

	_, err := usersCollection.UpdateOne(ctx, filter, update, setUpsert)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func (r *MongoReportsRepository) UpsertMonthlyReportAverageProfileViewTime(ctx context.Context, tweetId string, year int64, month int64, averageViewTime int64) error {
	_, span := r.tracer.Start(ctx, "MongoReportsRepository.UpsertMonthlyReportAverageProfileViewTime")
	defer span.End()

	usersCollection := r.cli.Database("reportsDB").Collection("reports")

	filter := bson.M{"tweetId": tweetId, "type": MONTHLY, "year": year, "month": month}
	update := bson.D{{"$set", bson.D{{"averageViewTime", averageViewTime}}}}
	setUpsert := options.Update().SetUpsert(true)

	_, err := usersCollection.UpdateOne(ctx, filter, update, setUpsert)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func (r *MongoReportsRepository) UpsertDailyReportLikesCount(ctx context.Context, tweetId string, year int64, month int64, day int64) error {
	_, span := r.tracer.Start(ctx, "MongoReportsRepository.UpsertDailyReportLikesCount")
	defer span.End()

	usersCollection := r.cli.Database("reportsDB").Collection("reports")

	filter := bson.M{"tweetId": tweetId, "type": DAILY, "year": year, "month": month, "day": day}
	update := bson.D{{"$inc", bson.D{{"likesCount", 1}}}}
	setUpsert := options.Update().SetUpsert(true)

	_, err := usersCollection.UpdateOne(ctx, filter, update, setUpsert)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func (r *MongoReportsRepository) UpsertDailyReportUnlikesCount(ctx context.Context, tweetId string, year int64, month int64, day int64) error {
	_, span := r.tracer.Start(ctx, "MongoReportsRepository.UpsertDailyReportUnlikesCount")
	defer span.End()

	usersCollection := r.cli.Database("reportsDB").Collection("reports")

	filter := bson.M{"tweetId": tweetId, "type": DAILY, "year": year, "month": month, "day": day}
	update := bson.D{{"$inc", bson.D{{"unlikesCount", 1}}}}
	setUpsert := options.Update().SetUpsert(true)

	_, err := usersCollection.UpdateOne(ctx, filter, update, setUpsert)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func (r *MongoReportsRepository) UpsertDailyReportProfileVisitsCount(ctx context.Context, tweetId string, year int64, month int64, day int64) error {
	_, span := r.tracer.Start(ctx, "MongoReportsRepository.UpsertDailyReportProfileVisitsCount")
	defer span.End()

	usersCollection := r.cli.Database("reportsDB").Collection("reports")

	filter := bson.M{"tweetId": tweetId, "type": DAILY, "year": year, "month": month, "day": day}
	update := bson.D{{"$inc", bson.D{{"profileVisits", 1}}}}
	setUpsert := options.Update().SetUpsert(true)

	_, err := usersCollection.UpdateOne(ctx, filter, update, setUpsert)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func (r *MongoReportsRepository) UpsertDailyReportAverageProfileViewTime(ctx context.Context, tweetId string, year int64, month int64, day int64, averageViewTime int64) error {
	_, span := r.tracer.Start(ctx, "MongoReportsRepository.UpsertDailyReportAverageProfileViewTime")
	defer span.End()

	usersCollection := r.cli.Database("reportsDB").Collection("reports")

	filter := bson.M{"tweetId": tweetId, "type": DAILY, "year": year, "month": month, "day": day}
	update := bson.D{{"$set", bson.D{{"averageViewTime", averageViewTime}}}}
	setUpsert := options.Update().SetUpsert(true)

	_, err := usersCollection.UpdateOne(ctx, filter, update, setUpsert)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
