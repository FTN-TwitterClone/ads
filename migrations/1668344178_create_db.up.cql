CREATE TABLE user_ads(
    username text,
    tweet_id uuid,
    timestamp timestamp,
    PRIMARY KEY ((username, tweet_id), timestamp)
) WITH CLUSTERING ORDER BY (timestamp DESC);
