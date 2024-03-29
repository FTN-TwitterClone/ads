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

type CassandraEventsRepository struct {
	tracer  trace.Tracer
	session *gocql.Session
}

func NewCassandraEventsRepository(tracer trace.Tracer) (*CassandraEventsRepository, error) {
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

	return &CassandraEventsRepository{
		tracer:  tracer,
		session: session,
	}, nil
}

func initKeyspace() error {
	dbport := os.Getenv("CASSANDRA_DBPORT")
	db := os.Getenv("CASSANDRA_DB")
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
	dbport := os.Getenv("CASSANDRA_DBPORT")
	db := os.Getenv("CASSANDRA_DB")
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

func (r *CassandraEventsRepository) SaveAdInfo(ctx context.Context, adInfo *model.AdInfo) error {
	_, span := r.tracer.Start(ctx, "CassandraEventsRepository.SaveTweetLikedEvent")
	defer span.End()

	err := r.session.Query("INSERT INTO ad_info(tweet_id, posted_by, town, min_age, max_age, gender) VALUES (?, ?, ?, ?, ?, ?)").
		Bind(adInfo.TweetId, adInfo.PostedBy, adInfo.Town, adInfo.MinAge, adInfo.MaxAge, adInfo.Gender).
		Exec()

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func (r *CassandraEventsRepository) GetAdInfo(ctx context.Context, tweetId string) (*model.AdInfo, error) {
	_, span := r.tracer.Start(ctx, "CassandraEventsRepository.GetAdInfo")
	defer span.End()

	var adInfo model.AdInfo

	err := r.session.Query("SELECT tweet_id, posted_by, town, min_age, max_age, gender FROM ad_info WHERE tweet_id = ?").
		Bind(tweetId).
		Scan(&adInfo.TweetId, &adInfo.PostedBy, &adInfo.Town, &adInfo.MinAge, &adInfo.MaxAge, &adInfo.Gender)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &adInfo, nil
}

func (r *CassandraEventsRepository) SaveTweetLikedEvent(ctx context.Context, tweetLikedEvent *model.TweetLikedEvent) error {
	_, span := r.tracer.Start(ctx, "CassandraEventsRepository.SaveTweetLikedEvent")
	defer span.End()

	err := r.session.Query("INSERT INTO tweet_liked_events(tweet_id, id, username) VALUES (?, ?, ?)").
		Bind(tweetLikedEvent.TweetId, gocql.UUIDFromTime(tweetLikedEvent.Time.UTC()), tweetLikedEvent.Username).
		Exec()

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func (r *CassandraEventsRepository) SaveTweetUnlikedEvent(ctx context.Context, tweetUnlikedEvent *model.TweetUnlikedEvent) error {
	_, span := r.tracer.Start(ctx, "CassandraEventsRepository.SaveTweetLikedEvent")
	defer span.End()

	err := r.session.Query("INSERT INTO tweet_unliked_events(tweet_id, id, username) VALUES (?, ?, ?)").
		Bind(tweetUnlikedEvent.TweetId, gocql.UUIDFromTime(tweetUnlikedEvent.Time.UTC()), tweetUnlikedEvent.Username).
		Exec()

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func (r *CassandraEventsRepository) SaveTweetViewedEvent(ctx context.Context, tweetViewedEvent *model.TweetViewedEvent) error {
	_, span := r.tracer.Start(ctx, "CassandraEventsRepository.SaveTweetLikedEvent")
	defer span.End()

	err := r.session.Query("INSERT INTO tweet_viewed_events(tweet_id, id, username, view_time) VALUES (?, ?, ?, ?)").
		Bind(tweetViewedEvent.TweetId, gocql.UUIDFromTime(tweetViewedEvent.Time.UTC()), tweetViewedEvent.Username, tweetViewedEvent.ViewTime).
		Exec()

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func (r *CassandraEventsRepository) SaveProfileVisitedEvent(ctx context.Context, profileVisitedEvent *model.ProfileVisitedEvent) error {
	_, span := r.tracer.Start(ctx, "CassandraEventsRepository.SaveProfileVisitedEvent")
	defer span.End()

	err := r.session.Query("INSERT INTO profile_visited_events(tweet_id, id, username) VALUES (?, ?, ?)").
		Bind(profileVisitedEvent.TweetId, gocql.UUIDFromTime(profileVisitedEvent.Time.UTC()), profileVisitedEvent.Username).
		Exec()

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func (r *CassandraEventsRepository) GetAverageTweetViewTime(ctx context.Context, tweetId gocql.UUID, from time.Time, to time.Time) (int, error) {
	_, span := r.tracer.Start(ctx, "CassandraAdsRepository.GetAverageTweetViewTime")
	defer span.End()

	var viewTime int

	err := r.session.Query("SELECT AVG(view_time) FROM tweet_viewed_events WHERE tweet_id = ? AND id > maxTimeuuid(?) AND id < minTimeuuid(?)").
		Bind(tweetId, from.UTC(), to.UTC()).
		Scan(&viewTime)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}

	return viewTime, err
}
