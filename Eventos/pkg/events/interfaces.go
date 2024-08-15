package events

import (
	"time"
)

type EventInterface interface {
	GetName() string
	GetDateTime() time.Time
	GetPayload() interface{}
}

type EventHandlerInterface interface {
	Handle(event EventInterface)
}

type EventDispatcherInterface interface {
	Register(eventName string, handler EventHandlerInterface)
	Dispatch(event EventInterface)
	Remove(eventName string, handler EventHandlerInterface)
	Has(eventName string, handler EventHandlerInterface) bool
	Clear()
}
