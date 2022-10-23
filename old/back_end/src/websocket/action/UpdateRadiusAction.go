package action

// UpdateRadiusAction is a request from client for update the article change listing radius
type UpdateRadiusAction struct {
	Action    ActionType `json:"action" form:"action"  binding:"required"`
	NewRadius float32    `form:"new_radius"  binding:"required" json:"new_radius"`
}

func (a UpdateRadiusAction) GetActionType() ActionType {
	return a.Action
}
