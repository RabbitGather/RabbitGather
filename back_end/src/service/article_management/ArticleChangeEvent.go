package article_management

import "encoding/json"

type ArticleChangeEvent struct {
	ID        int64          `json:"id"`
	Event     string         `json:"event"`
	Timestamp int64          `json:"timestamp"`
	Position  PositionStruct `json:"position"`
}

func (e *ArticleChangeEvent) ToJsonString() string {
	s, err := json.Marshal(*e)
	if err != nil {
		log.ERROR.Println("error when marshal ArticleChangeEvent")
	}
	return string(s)
}
