package events

import "rabbit_gather/src/service/article_management"

type ArticleChangeEvent struct {
	Event     string                            `json:"event"`
	ID        int64                             `json:"id ,omitempty"`
	Timestamp int64                             `json:"timestamp ,omitempty"`
	Position  article_management.PositionStruct `json:"position ,omitempty"`
}
