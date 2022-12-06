package cassandra

import (
	"context"
	"fmt"
	"github.com/FTN-TwitterClone/ads/model"
	"github.com/gocql/gocql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/cassandra"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"log"
	"os"
	"time"
)

type CassandraAdsRepository struct {
	tracer  trace.Tracer
	session *gocql.Session
}

func NewCassandraAdsRepository(tracer trace.Tracer) (*CassandraAdsRepository, error) {
	err := initKeyspace()
	if err != nil {
		return nil, err
	}

	err = migrateDB()
	if err != nil {
		return nil, err
	}

	dbport := os.Getenv("CASSANDRA_DBPORT")
	db := os.Getenv("CASSANDRA_DB")
	host := fmt.Sprintf("%s:%s", db, dbport)

	cluster := gocql.NewCluster(host)
	cluster.ProtoVersion = 4
	cluster.Keyspace = "ads_database"
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	//defer session.Close()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Printf("Connected OK!")

	return &CassandraAdsRepository{
		tracer:  tracer,
		session: session,
	}, nil
}

func initKeyspace() error {
	dbport := os.Getenv("DBPORT")
	db := os.Getenv("DB")
	host := fmt.Sprintf("%s:%s", db, dbport)

	cluster := gocql.NewCluster(host)
	cluster.ProtoVersion = 4
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	defer session.Close()

	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Printf("Connected OK!")

	err = session.Query("CREATE KEYSPACE IF NOT EXISTS ads_database WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1}").Exec()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func migrateDB() error {
	dbport := os.Getenv("DBPORT")
	db := os.Getenv("DB")
	connString := fmt.Sprintf("cassandra://%s:%s/ads_database?x-multi-statement=true", db, dbport)

	m, err := migrate.New("file://migrations", connString)
	if err != nil {
		log.Fatal(err)
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
		return err
	}

	return nil
}

func (r *CassandraAdsRepository) SaveTweetLikedEvent(ctx context.Context, tweetLikedEvent *model.TweetLikedEvent) error {
	_, span := r.tracer.Start(ctx, "CassandraAdsRepository.SaveTweetLikedEvent")
	defer span.End()

	err := r.session.Query("INSERT INTO tweet_liked_events(tweet_id, id, username) VALUES (?, ?, ?)").
		Bind(tweetLikedEvent.TweetId, gocql.TimeUUID(), tweetLikedEvent.Username).
		Exec()

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func (r *CassandraAdsRepository) SaveTweetUnlikedEvent(ctx context.Context, tweetUnlikedEvent *model.TweetUnlikedEvent) error {
	_, span := r.tracer.Start(ctx, "CassandraAdsRepository.SaveTweetLikedEvent")
	defer span.End()

	err := r.session.Query("INSERT INTO tweet_unliked_events(tweet_id, id, username) VALUES (?, ?, ?)").
		Bind(tweetUnlikedEvent.TweetId, gocql.TimeUUID(), tweetUnlikedEvent.Username).
		Exec()

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func (r *CassandraAdsRepository) SaveTweetViewedEvent(ctx context.Context, tweetViewedEvent *model.TweetViewedEvent) error {
	_, span := r.tracer.Start(ctx, "CassandraAdsRepository.SaveTweetLikedEvent")
	defer span.End()

	err := r.session.Query("INSERT INTO tweet_viewed_events(tweet_id, id, username, view_time) VALUES (?, ?, ?, ?)").
		Bind(tweetViewedEvent.TweetId, gocql.TimeUUID(), tweetViewedEvent.Username, tweetViewedEvent.ViewTime).
		Exec()

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func (r *CassandraAdsRepository) SaveProfileVisitedEvent(ctx context.Context, profileVisitedEvent *model.ProfileVisitedEvent) error {
	_, span := r.tracer.Start(ctx, "CassandraAdsRepository.SaveProfileVisitedEvent")
	defer span.End()

	err := r.session.Query("INSERT INTO profile_visited_events(tweet_id, id, username) VALUES (?, ?, ?)").
		Bind(profileVisitedEvent.TweetId, gocql.TimeUUID(), profileVisitedEvent.Username).
		Exec()

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func (r *CassandraAdsRepository) GetTweetLikesCount(ctx context.Context, tweetId gocql.UUID, from time.Time, to time.Time) (int, error) {
	_, span := r.tracer.Start(ctx, "CassandraAdsRepository.GetProfileVisitsCount")
	defer span.End()

	var visitsCount int

	err := r.session.Query("SELECT COUNT(*) FROM tweet_liked_events WHERE tweet_id = ? AND id > maxTimeuuid(?) AND id < minTimeuuid(?)").
		Bind(tweetId, from, to).
		Scan(&visitsCount)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}

	return visitsCount, err
}

func (r *CassandraAdsRepository) GetTweetUnlikesCount(ctx context.Context, tweetId gocql.UUID, from time.Time, to time.Time) (int, error) {
	_, span := r.tracer.Start(ctx, "CassandraAdsRepository.GetProfileVisitsCount")
	defer span.End()

	var visitsCount int

	err := r.session.Query("SELECT COUNT(*) FROM tweet_unliked_events WHERE tweet_id = ? AND id > maxTimeuuid(?) AND id < minTimeuuid(?)").
		Bind(tweetId, from, to).
		Scan(&visitsCount)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}

	return visitsCount, err
}

func (r *CassandraAdsRepository) GetAverageTweetViewTimeCount(ctx context.Context, tweetId gocql.UUID, from time.Time, to time.Time) (int, error) {
	_, span := r.tracer.Start(ctx, "CassandraAdsRepository.GetProfileVisitsCount")
	defer span.End()

	var visitsCount int

	err := r.session.Query("SELECT AVG(view_time) FROM tweet_viewed_events WHERE tweet_id = ? AND id > maxTimeuuid(?) AND id < minTimeuuid(?)").
		Bind(tweetId, from, to).
		Scan(&visitsCount)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}

	return visitsCount, err
}

func (r *CassandraAdsRepository) GetProfileVisitsCount(ctx context.Context, tweetId gocql.UUID, from time.Time, to time.Time) (int, error) {
	_, span := r.tracer.Start(ctx, "CassandraAdsRepository.GetProfileVisitsCount")
	defer span.End()

	var visitsCount int

	err := r.session.Query("SELECT COUNT(*) FROM profile_visited_events WHERE tweet_id = ? AND id > maxTimeuuid(?) AND id < minTimeuuid(?)").
		Bind(tweetId, from, to).
		Scan(&visitsCount)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}

	return visitsCount, err
}
