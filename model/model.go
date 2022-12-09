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
	TweetId  gocql.UUID `json:"tweetId"`
	PostedBy string     `json:"postedBy"`
	Town     string     `json:"town"`
	MinAge   int32      `json:"minAge"`
	MaxAge   int32      `json:"maxAge"`
	Gender   string     `json:"gender"`
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
	ViewTime int32
	Time     time.Time
}

type ProfileVisitedEvent struct {
	Username string
	TweetId  gocql.UUID
	Time     time.Time
}

type TweetViewTime struct {
	ViewTime int32 `json:"viewTime"`
}

type Report struct {
	TweetId         string `json:"tweetId" bson:"tweetId"`
	Year            int64  `json:"year" bson:"year"`
	Month           int64  `json:"month" bson:"month"`
	Day             int64  `json:"day" bson:"day"`
	LikesCount      int    `json:"tweetsLiked" bson:"tweetsLiked"`
	UnlikesCount    int    `json:"tweetsUnliked" bson:"tweetsUnliked"`
	ProfileVisits   int    `json:"profileVisits" bson:"profileVisits"`
	AverageViewTime int    `json:"averageViewTime" bson:"averageViewTime"`
}
