CREATE TABLE tweet_liked_events(
    tweet_id timeuuid,
    id timeuuid,
    username text,
    PRIMARY KEY ((tweet_id), id)
);

CREATE TABLE tweet_unliked_events(
    tweet_id timeuuid,
    id timeuuid,
    username text,
    PRIMARY KEY ((tweet_id), id)
);

CREATE TABLE tweet_viewed_events(
    tweet_id timeuuid,
    id timeuuid,
    username text,
    view_time int,
    PRIMARY KEY ((tweet_id), id)
);

CREATE TABLE profile_visited_events(
    tweet_id timeuuid,
    id timeuuid,
    username text,
    PRIMARY KEY ((tweet_id), id)
);

CREATE TABLE reports(
    id uuid PRIMARY KEY,
    likes_added INT,
    likes_removed INT,
    avg_read_time INT,
    num_of_profile_visits INT,
    from_date timestamp,
    to_date timestamp
);
