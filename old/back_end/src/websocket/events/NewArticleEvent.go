package events

import (
	"rabbit_gather/util"
)

// The NewArticleEvent is an event represented a new article been created.
type NewArticleEvent struct {
	Event     EventType    `json:"event" form:"event"  binding:"required"`
	ArticleID int64        `json:"article_id ,omitempty"`
	Position  util.Point2D `form:"position"  binding:"required" json:"position"`
	Timestamp int64        `json:"timestamp"`
}

func (a NewArticleEvent) GetEventType() EventType {
	return a.Event
}
