package databinding

import "errors"

const (
	ErrPropertyReadOnly = errors.New("read-only property")
)

type eventType uint32

const (
	EventValueWillChange eventType = 1 << iota
	EventValueDidChange
)

type Event interface {
	Type() eventType
}

type WatcherChan chan Event

type Property interface {
	Get() interface{}
	Set(interface{}) error
	Watch() WatcherChan
	Dispose()
}
