package protomanager

type EventType int

const (
	EventTypeSuccess EventType = iota
	EventTypeError
	EventTypeInfo
)

type Event struct {
	Type    EventType
	Message string
}

// SubscribeEvents returns a channel for receiving event updates.
func (pm *ProtoManager) SubscribeEvents() <-chan Event {
	eventCh := make(chan Event, 10)
	go func() {
		for {
			// Simulated event generation (in actual implementation, this would respond to real events)
			eventCh <- Event{Type: EventTypeInfo, Message: "Sample event"}
		}
	}()
	return eventCh
}