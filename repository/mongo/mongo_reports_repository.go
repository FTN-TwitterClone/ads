package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel/trace"
	"os"
)

type MongoReportsRepository struct {
	tracer trace.Tracer
	cli    *mongo.Client
}

func NewMongoReportsRepository(tracer trace.Tracer) (*MongoReportsRepository, error) {

	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")

	//mongo logic
	host := fmt.Sprintf("%s:%s", db, dbport)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(`mongodb://`+host))
	if err != nil {
		panic(err)
	}
	client.Database("reportsDB").Collection("users")

	car := MongoReportsRepository{
		tracer,
		client,
	}

	return &car, nil
}

func (r *MongoReportsRepository) UpsertMonthlyReportTweetLiked(ctx context.Context, tweetId string) error {
	return nil
}
