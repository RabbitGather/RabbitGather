package action

// ListenOnNewArticleAction is a request from client for subscribe a new article to listener.
type ListenOnNewArticleAction struct {
	Action    ActionType `json:"action" form:"action"  binding:"required"`
	ArticleID string     `form:"article_id"  binding:"required" json:"article_id"`
}

func (a ListenOnNewArticleAction) GetActionType() ActionType {
	return a.Action
}
