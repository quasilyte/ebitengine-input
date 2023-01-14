package input

// SimulatedEvent represents a virtual input that can be send down the stream.
//
// The data carried by this event will be used to construct an EventInfo object.
//
// Experimental: this is a part of virtual input API, which is not stable yet.
type SimulatedEvent struct {
	Key Key

	Pos Vec
}

// EventInfo holds extra information about the input device event.
//
// Pos carries the event location, if available.
// Pos is a click location for mouse events.
// Pos is a tap location for screen touch events.
// Use HasPos() predicate to know whether there is a pos associated
// with the event to distinguish between (0, 0) pos and lack of pos info.
type EventInfo struct {
	kind   keyKind
	hasPos bool

	Pos Vec
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
