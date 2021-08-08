package events

type ArticleErrorEvent struct {
	Event   string `json:"event"`
	Message string `json:"message"`
}
