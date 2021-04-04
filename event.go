package main

import "sync"

// can be EventInfo, EventWarn or EventErr
type EventType string

// EventType can now implement the Stringer interface
func (t EventType) String() string {
	return string(t)
}

// EventRegister stores a channel which accepts messages and a logger.
// It's purpose is to prevent calling 2 methods each time some event
// should be sent to the frontend and logged, instead one method is called,
// EventRegister.Register which will do both things.
type EventRegister struct {
	// makes sure that all logs are logged in right order
	mu      sync.Mutex
	logger  *Logger
	channel chan string
}

// NewEventRegister returns a pointer to a new EventRegister. Channel can be nil.
func NewEventRegister(logger *Logger, channel chan string) *EventRegister {
	r := EventRegister{sync.Mutex{}, logger, channel}
	return &r
}

// Register registers the event aka send a status about the event through the
// channel and logs it via Logger.
func (r *EventRegister) Register(evt EventType, logtext, chantext string) {
	if logtext != "" {
		r.mu.Lock()
		r.logger.Log(evt, logtext)
		r.mu.Unlock()
	}
	if chantext != "" && r.channel != nil {
		r.channel <- chantext
	}
}
