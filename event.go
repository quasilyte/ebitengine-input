package input

// SimulatedKeyEvent represents a virtual input that can be send down the stream.
//
// The data carried by this event will be used to construct an EventInfo object.
//
// Experimental: this is a part of virtual input API, which is not stable yet.
type SimulatedKeyEvent struct {
	Key Key

	Pos Vec
}

// SimulatedAction represents an artificially triggered action.
//
// It shares many properties with SimulatedKeyEvent, but
// the event consumers will have no way of knowing which input
// device was used to emit this event, because SimulatedAction
// has no device associated with it.
//
// As a consequence, all event info methods like IsGamepadeEvent() will report false.
// It's possible to trigger an action that has no keys associated with it.
// All actions triggered using this method will be only visible to the handler
// of the same player ID (like gamepad button events).
//
// Experimental: this is a part of virtual input API, which is not stable yet.
type SimulatedAction struct {
	Action Action

	Pos Vec

	StartPos Vec
}

// EventInfo holds extra information about the input device event.
//
// Pos carries the event location, if available.
// Pos is a click location for mouse events.
// Pos is a tap location for screen touch events.
// Use HasPos() predicate to know whether there is a pos associated
// with the event to distinguish between (0, 0) pos and lack of pos info.
//
// StartPos is only set for a few events where it makes sense.
// A drag event, for instance, will store the "dragging from" location there.
type EventInfo struct {
	kind   keyKind
	hasPos bool

	Pos      Vec
	StartPos Vec
}

// HasPos reports whether this event has a position associated with it.
// Use Pos field to get the pos value.
func (e EventInfo) HasPos() bool { return e.hasPos }

// IsTouchEvent reports whether this event was triggered by a screen touch device.
func (e EventInfo) IsTouchEvent() bool { return e.kind == keyTouch }

// IsKeyboardEvent reports whether this event was triggered by a keyboard device.
func (e EventInfo) IsKeyboardEvent() bool { return e.kind == keyKeyboard }

// IsMouseEvent reports whether this event was triggered by a mouse device.
func (e EventInfo) IsMouseEvent() bool { return e.kind == keyMouse }

// IsGamepadEvent reports whether this event was triggered by a gamepad device.
func (e EventInfo) IsGamepadEvent() bool {
	switch e.kind {
	case keyGamepad, keyGamepadLeftStick, keyGamepadRightStick:
		return true
	default:
		return false
	}
}

type simulatedEvent struct {
	code     int
	keyKind  keyKind
	playerID uint8

	pos      Vec
	startPos Vec
}
