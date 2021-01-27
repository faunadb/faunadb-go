package faunadb

import (
	"fmt"
)

// StreamEventType is a stream eveny type
type StreamEventType = string

const (
	// ErrorEventT is the stream error event type
	ErrorEventT StreamEventType = "error"

	// HistoryRewriteEventT is the stream history rewrite event type
	HistoryRewriteEventT StreamEventType = "history_rewrite"

	// StartEventT is the stream start event type
	StartEventT StreamEventType = "start"

	// VersionEventT is the stream version event type
	VersionEventT StreamEventType = "version"
)

// StreamEvent represents a stream event with a `type` and `txn`
type StreamEvent interface {
	Type() StreamEventType
	Txn() int64
	String() string
}

// StartEvent emitted when a valid stream subscription begins.
// Upcoming events are guaranteed to have transaction timestamps equal to or greater than
// the stream's start timestamp.
type StartEvent struct {
	StreamEvent
	txn   int64
	event int64
}

func unMarshalStreamEvent(data Obj) (evt StreamEvent, err error) {
	switch StreamEventType(data["type"].(StringV)) {
	case StartEventT:
		evt = StartEvent{
			txn:   int64(data["txn"].(LongV)),
			event: int64(data["event"].(LongV)),
		}
	case VersionEventT:
		evt = VersionEvent{
			txn:   int64(data["txn"].(LongV)),
			event: data["event"].(ObjectV),
		}
	case ErrorEventT:
		evt = ErrorEvent{
			txn: int64(data["txn"].(LongV)),
			err: errorFromStreamError(data["event"].(ObjectV)),
		}

	case HistoryRewriteEventT:
		evt = HistoryRewriteEvent{
			txn:   int64(data["txn"].(LongV)),
			event: data["event"].(ObjectV),
		}
	}
	return
}

// Type returns the stream event type
func (event StartEvent) Type() StreamEventType {
	return StartEventT
}

// Txn returns the stream event timestamp
func (event StartEvent) Txn() int64 {
	return event.txn
}

// Event returns the stream event as a `f.ObjectV`
func (event StartEvent) Event() int64 {
	return event.event
}

func (event StartEvent) String() string {
	return fmt.Sprintf("StartEvent{event=%d, txn=%d} ", event.Event(), event.Txn())
}

// VersionEvent represents a version event that occurs upon any
// modifications to the current state of the subscribed document.
type VersionEvent struct {
	StreamEvent
	txn   int64
	event ObjectV
}

// Txn returns the stream event timestamp
func (event VersionEvent) Txn() int64 {
	return event.txn
}

// Event returns the stream event as a `f.ObjectV`
func (event VersionEvent) Event() ObjectV {
	return event.event
}

func (event VersionEvent) String() string {
	return fmt.Sprintf("VersionEvent{txn=%d, event=%s}", event.Txn(), event.Event())
}

// Type returns the stream event type
func (event VersionEvent) Type() StreamEventType {
	return VersionEventT
}

// HistoryRewriteEvent represents a history rewrite event which occurs upon any modifications
// to the history of the subscribed document.
type HistoryRewriteEvent struct {
	StreamEvent
	txn   int64
	event ObjectV
}

// Txn returns the stream event timestamp
func (event HistoryRewriteEvent) Txn() int64 {
	return event.txn
}

// Event returns the stream event as a `f.ObjectV`
func (event HistoryRewriteEvent) Event() ObjectV {
	return event.event
}

func (event HistoryRewriteEvent) String() string {
	return fmt.Sprintf("HistoryRewriteEvent{txn=%d, event=%s}", event.Txn(), event.Event())
}

// Type returns the stream event type
func (event HistoryRewriteEvent) Type() StreamEventType {
	return VersionEventT
}

// ErrorEvent represents an error event fired both for client and server errors
// that may occur as a result of a subscription.
type ErrorEvent struct {
	StreamEvent
	txn int64
	err error
}

// Type returns the stream event type
func (event ErrorEvent) Type() StreamEventType {
	return ErrorEventT
}

// Txn returns the stream event timestamp
func (event ErrorEvent) Txn() int64 {
	return event.txn
}

// Error returns the event error
func (event ErrorEvent) Error() error {
	return event.err
}

func (event ErrorEvent) String() string {
	return fmt.Sprintf("ErrorEvent{error=%s}", event.err)
}