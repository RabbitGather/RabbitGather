package events

// The DeleteArticleEvent is an event represented an article been deleted
type DeleteArticleEvent struct {
	Event     EventType `json:"event" form:"event"  binding:"required"`
	ArticleID int64     `json:"article_id ,omitempty"`
}

func (a DeleteArticleEvent) GetEventType() EventType {
	return a.Event
}
