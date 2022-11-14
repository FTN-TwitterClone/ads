package model

import (
	"github.com/gocql/gocql"
	"time"
)

// Info from JWT token
type AuthUser struct {
	Username string
	Role     string
	Exp      time.Time
}

type TweetLikedEvent struct {
	Username string
	TweetId  gocql.UUID
	Time     time.Time
}

type TweetUnlikedEvent struct {
	Username string
	TweetId  gocql.UUID
	Time     time.Time
}

type TweetViewedEvent struct {
	Username string
	TweetId  gocql.UUID
	ReadTime int32
	Time     time.Time
}

type ProfileVisitedEvent struct {
	Username string
	TweetId  gocql.UUID
	Time     time.Time
}

type Report struct {
	Id                   gocql.UUID `json:"id"`
	TweetId              gocql.UUID `json:"tweetId"`
	From                 time.Time  `json:"from"`
	To                   time.Time  `json:"to"`
	TimeGenerated        time.Time  `json:"timeGenerated"`
	TweetsLiked          int        `json:"tweetsLiked"`
	TweetsUnliked        int        `json:"tweetsUnliked"`
	AverageTweetViewTime int        `json:"averageTweetViewTime"`
	ProfileVisits        int        `json:"profileVisits"`
}
