package action

type ActionType string

const (
	UpdateRadius       ActionType = "update_radius_action"
	TextMessage        ActionType = "text_action"
	ListenOnNewArticle ActionType = "listen_on_new_article"
)
const ACTION_KEY = "action"

// The Action is the data packet came from client
type Action interface {
	GetActionType() ActionType
}
