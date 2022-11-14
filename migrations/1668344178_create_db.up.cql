CREATE TABLE likesAdded(
    id uuid,
    tweetId uuid,
    username text,
    time timestamp,
    PRIMARY KEY ((id), time)
);

CREATE INDEX ON likesAdded (tweetid);

CREATE TABLE likesRemoved(
    id uuid,
    tweetId uuid,
    username text,
    time timestamp,
    PRIMARY KEY ((id), time)
);

CREATE INDEX ON likesRemoved (tweetid);

CREATE TABLE tweetReadTimes(
    id uuid,
    tweetId uuid,
    username text,
    readTime INT,
    time timestamp,
    PRIMARY KEY ((id), time)
);

CREATE INDEX ON tweetReadTimes (tweetid);

CREATE TABLE profileVisits(
    id uuid,
    tweetId uuid,
    username text,
    time timestamp,
    PRIMARY KEY ((id), time)
);

CREATE INDEX ON profileVisits (tweetid);

CREATE TABLE reports(
    id uuid PRIMARY KEY,
    likesAdded INT,
    likesRemoved INT,
    avgReadTime INT,
    numOfProfileVisits INT,
    fromDate timestamp,
    toDate timestamp
);
