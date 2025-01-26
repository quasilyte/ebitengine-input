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
//
// Duration carries the key press duration if available.
// Duration specifies how long the key has been pressed in ticks same as inpututil.KeyPressDuration.
// Duration for key press with modifiers it will return the lowest duration of all key presses.
// Use HasDuration() predicate to know whether there is a duration associated
// with the event to distinguish between 0 duration and lack of duration info.
type EventInfo struct {
	kind        keyKind
	hasPos      bool
	hasDuration bool

	Duration int
	Pos      Vec
	StartPos Vec
}

// HasPos reports whether this event has a position associated with it.
// Use Pos field to get the pos value.
func (e EventInfo) HasPos() bool { return e.hasPos }

// HasDuration reports whether this event has a press duration associated with it.
// Use Duration field to get the press duration value.
func (e EventInfo) HasDuration() bool { return e.hasDuration }

// IsTouchEvent reports whether this event was triggered by a screen touch device.
//
// Deprecated: Use Source().IsTouch() instead.
func (e EventInfo) IsTouchEvent() bool {
	return e.Source().IsTouch()
}

// IsKeyboardEvent reports whether this event was triggered by a keyboard device.
//
// Deprecated: Use Source().IsKeyboard() instead.
func (e EventInfo) IsKeyboardEvent() bool {
	return e.Source().IsKeyboard()
}

// IsMouseEvent reports whether this event was triggered by a mouse device.
//
// Deprecated: Use Source().IsMouse() instead.
func (e EventInfo) IsMouseEvent() bool {
	return e.Source().IsMouse()
}

// IsGamepadEvent reports whether this event was triggered by a gamepad device.
//
// Deprecated: Use Source().IsGamepad() instead.
func (e EventInfo) IsGamepadEvent() bool {
	return e.Source().IsGamepad()
}

// Source returns the set of devices that were used to trigger the event.
//
// Usually, it returns a single device mask, but sometimes an action
// can be triggered by several devices.
// For example, a ctrl+click key would involve both keyboard and mouse devices,
// resulting in a KeyboardDevice|KeyboardMouse value.
//
// A simulated event would have a zero mask returned, meaning
// no real device was involved.
func (e EventInfo) Source() DeviceKind {
	return e.kind.device()
}

type simulatedEvent struct {
	code     int
	keyKind  keyKind
	playerID uint8

	pos      Vec
	startPos Vec
}
