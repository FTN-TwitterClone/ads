CREATE TABLE ad_info(
    tweet_id timeuuid,
    posted_by text,
    town text,
    min_age int,
    max_age int,
    gender text,
    PRIMARY KEY ((tweet_id))
);
