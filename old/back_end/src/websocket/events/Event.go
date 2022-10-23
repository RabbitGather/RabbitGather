package events

type EventType string

const (
	NewArticle    EventType = "new_article"
	DeleteArticle EventType = "delete_article"
)

// The Event is the data packet that are going to send to Client
type Event interface {
	GetEventType() EventType
}
