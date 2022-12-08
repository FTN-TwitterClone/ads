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

type AdInfo struct {
	TweetId  gocql.UUID
	PostedBy string
	Town     string
	MinAge   int32
	MaxAge   int32
	Gender   string
}

type TweetLikedEvent struct {
	Username string
	TweetId  gocql.UUID
}

type TweetUnlikedEvent struct {
	Username string
	TweetId  gocql.UUID
}

type TweetViewedEvent struct {
	Username string
	TweetId  gocql.UUID
	ViewTime int32
}

type ProfileVisitedEvent struct {
	Username string
	TweetId  gocql.UUID
}

type TweetViewTime struct {
	ViewTime int32 `json:"viewTime"`
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
