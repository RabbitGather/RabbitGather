package events

import (
	"rabbit_gather/util"
)

// The ArticleChangeEvent is an event represented an article changed
type ArticleChangeEvent struct {
	Event     EventType    `json:"event" form:"event"  binding:"required"`
	ID        int64        `json:"id ,omitempty"`
	Timestamp int64        `json:"timestamp ,omitempty"`
	Position  util.Point2D `json:"position ,omitempty"`
}

func (a ArticleChangeEvent) GetEventType() EventType {
	return a.Event
}
