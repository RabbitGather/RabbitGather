package events

// The ErrorEvent is an event represented an error occurred.
type ErrorEvent struct {
	Event   EventType `json:"event" form:"event"  binding:"required"`
	Message string    `json:"message"`
}

func (a ErrorEvent) GetEventType() EventType {
	return a.Event
}
