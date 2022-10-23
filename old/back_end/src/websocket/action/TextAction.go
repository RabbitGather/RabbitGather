package action

// TextAction is a request structure of plain text message from client.
type TextAction struct {
	Action ActionType `json:"action" form:"action"  binding:"required"`
	Text   string     `form:"text"  binding:"required" json:"text"`
}

func (a TextAction) GetActionType() ActionType {
	return a.Action
}
