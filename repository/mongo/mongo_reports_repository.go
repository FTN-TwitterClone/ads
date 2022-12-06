package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"os"
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

func (r *MongoReportsRepository) UpsertMonthlyReportLikesCount(ctx context.Context, tweetId string, year int64, month int64) error {
	_, span := r.tracer.Start(ctx, "MongoReportsRepository.UpsertMonthlyReportLikesCount")
	defer span.End()

	usersCollection := r.cli.Database("reportsDB").Collection("reports")

	filter := bson.M{"tweetId": tweetId, "type": "monthly", "year": year, "month": month}
	update := bson.D{{"inc", bson.D{{"likesCount", 1}}}}
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

	filter := bson.M{"tweetId": tweetId, "type": "monthly", "year": year, "month": month}
	update := bson.D{{"set", bson.D{{"averageViewTime", averageViewTime}}}}
	setUpsert := options.Update().SetUpsert(true)

	_, err := usersCollection.UpdateOne(ctx, filter, update, setUpsert)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
