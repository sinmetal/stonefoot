package main

import (
	"context"
	"net/http"

	"github.com/favclip/ucon"
	"github.com/favclip/ucon/swagger"
)

func SetUpTweetAPI(swPlugin *swagger.Plugin) {
	api := &TweetAPI{}
	tag := swPlugin.AddTag(&swagger.Tag{Name: "Tweet", Description: "Tweet API list"})
	var hInfo *swagger.HandlerInfo

	hInfo = swagger.NewHandlerInfo(api.Insert)
	ucon.Handle(http.MethodPost, "/api/1/tweet", hInfo)
	hInfo.Description, hInfo.Tags = "post to tweet", []string{tag.Name}
}

// TweetAPI is Tweetに関するAPI
type TweetAPI struct{}

type TweetAPIPostRequest struct {
	ID      string `json:"id",swagger:"req"`
	Author  string `json:"author"`
	Content string `json:"content"`
	Sort    int64  `json:"sort"`
}

// TweetAPIPostResponse is TweetAPI PostのResponse
type TweetAPIPostResponse struct {
	Entity *TweetV2 `json:"entity"`
}

func (api *TweetAPI) Insert(ctx context.Context, req *TweetAPIPostRequest) (*TweetAPIPostResponse, error) {
	s, err := NewTweetV2Store(ctx, spannerClient)
	if err != nil {
		return nil, err
	}
	t, err := s.Insert(ctx, &TweetV2{
		ID:      req.ID,
		Author:  req.Author,
		Content: req.Content,
		Favos:   []string{},
		Sort:    req.Sort,
	})
	if err != nil {
		return nil, err
	}
	return &TweetAPIPostResponse{Entity: t}, nil
}
