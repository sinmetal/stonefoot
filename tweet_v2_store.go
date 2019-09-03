package main

import (
	"context"
	"hash/crc32"
	"time"

	"cloud.google.com/go/spanner"
)

type TweetV2 struct {
	ID             string `spanner:"id"`
	Author         string
	Content        string
	Count          int64
	Favos          []string
	Sort           int64
	ShardCreatedAt int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CommitedAt     time.Time
}

type TweetV2Store struct {
	sc *spanner.Client
}

func NewTweetV2Store(ctx context.Context, sc *spanner.Client) (*TweetV2Store, error) {
	return &TweetV2Store{
		sc: sc,
	}, nil
}

func (s *TweetV2Store) TableName() string {
	return "TweetV2"
}

func (s *TweetV2Store) Insert(ctx context.Context, tweet *TweetV2) (*TweetV2, error) {
	ctx, span := StartSpan(ctx, "/tweet/insert")
	defer span.End()

	now := time.Now()
	tweet.ShardCreatedAt = int64(crc32.ChecksumIEEE([]byte(now.String())) % 10)
	tweet.CreatedAt = now
	tweet.UpdatedAt = now
	tweet.CommitedAt = spanner.CommitTimestamp

	m, err := spanner.InsertStruct(s.TableName(), tweet)
	if err != nil {
		return nil, err
	}
	ms := []*spanner.Mutation{
		m,
	}

	commit, err := s.sc.Apply(ctx, ms)
	if err != nil {
		return nil, err
	}
	tweet.CommitedAt = commit

	return tweet, nil
}
