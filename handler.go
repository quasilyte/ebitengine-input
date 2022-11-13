package input

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Handler is used to associate a keymap with an abstract input consumer.
//
// The ID bound to the handler is used to distinguish which gamepad is
// related to this handler.
//
// You usually need to create the input handlers only once and carry
// them through the game using your preferred method.
//
// If any game object needs to handle the input, they need an input handler object.
type Handler struct {
	id     int
	keymap Keymap
	sys    *System
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

// GamepadConnected reports whether the gamepad associated with this handler is connected.
// The gamepad ID is the handler ID used during the handler creation.
//
// There should be at least one call to the System.Update() before this function
// can return the correct results.
func (h *Handler) GamepadConnected() bool {
	for _, id := range h.sys.gamepadIDs {
		if id == ebiten.GamepadID(h.id) {
			return true
		}
	}
	return false
}

// TouchEventsEnabled reports whether this handler can receive screen touch events.
func (h *Handler) TouchEventsEnabled() bool {
	return h.sys.touchEnabled
}

// TapPos is like CursorPos(), but for the screen tapping.
// If there is no screen tapping in this frame, it returns false.
func (h *Handler) TapPos() (Vec, bool) {
	return h.sys.touchTapPos, h.sys.touchHasTap
}

// CursorPos returns the current mouse cursor position on the screen.
func (h *Handler) CursorPos() Vec {
	return h.sys.cursorPos
}

// DefaultInputMask returns the input mask suitable for functions like ActionKeyNames.
//
// If gamepad is connected, it returns GamepadInput mask.
// Otherwise it returns KeyboardInput+MouseInput mask.
// This is good enough for the simplest games, but you may to implement this
// logic inside your game if you need something more complicated.
func (h *Handler) DefaultInputMask() InputDeviceKind {
	if h.GamepadConnected() {
		return GamepadInput
	}
	return KeyboardInput | MouseInput
}

// ActionKeyNames returns a list of key names associated by this action.
//
// It filters the results by a given input device mask.
// If you want to include all input device keys, use AnyInput value.
//
// This function is useful when you want to display a list of keys
// the player should press in order to activate some action.
//
// The filtering is useful to avoid listing the unrelated options.
// For example, if player uses the gamepad, it could be weird to
// show keyboard options listed. For the simple cases, you can use
// DefaultInputMask() method to get the mask that will try to avoid
// that situation. See its comment to learn more.
func (h *Handler) ActionKeyNames(action Action, mask InputDeviceKind) []string {
	keys, ok := h.keymap[action]
	if !ok {
		return nil
	}
	gamepadConnected := h.GamepadConnected()
	result := make([]string, 0, len(keys))
	for _, k := range keys {
		enabled := true
		switch k.kind {
		case keyKeyboard:
			enabled = mask&KeyboardInput != 0
		case keyMouse:
			enabled = mask&MouseInput != 0
		case keyGamepad, keyGamepadLeftStick, keyGamepadRightStick:
			enabled = gamepadConnected && (mask&GamepadInput != 0)
		case keyTouch:
			enabled = h.sys.touchEnabled && (mask&TouchInput != 0)
		}
		if enabled {
			result = append(result, k.name)
		}
	}
	return result
}

// JustPressedActionInfo is like ActionIsJustPressed, but with more information.
//
// The second return value is false is given action is not activated.
//
// The first return value will hold the extra event info.
// See EventInfo comment to learn more.
func (h *Handler) JustPressedActionInfo(action Action) (EventInfo, bool) {
	var info EventInfo
	keys, ok := h.keymap[action]
	if !ok {
		return info, false
	}
	isPressed := false
	for _, k := range keys {
		if !h.keyIsJustPressed(k) {
			continue
		}
		isPressed = true
		info.kind = k.kind
		switch k.kind {
		case keyMouse:
			info.Pos = h.sys.cursorPos
			info.hasPos = true
			return info, true
		case keyTouch:
			info.Pos = h.sys.touchTapPos
			info.hasPos = true
			return info, true
		}
	}
	return info, isPressed
}

// ActionIsJustPressed is like ebitenutil.IsKeyJustPressed, but operates
// on the action level and works with any kinds of "keys".
// It returns true if any of the keys bound to the action was pressed during this frame.
func (h *Handler) ActionIsJustPressed(action Action) bool {
	keys, ok := h.keymap[action]
	if !ok {
		return false
	}
	for _, k := range keys {
		if h.keyIsJustPressed(k) {
			return true
		}
	}
	return false
}

// ActionIsPressed is like ebiten.IsKeyPressed, but operates
// on the action level and works with any kinds of "keys".
// It returns true if any of the keys bound to the action is being pressed.
func (h *Handler) ActionIsPressed(action Action) bool {
	keys, ok := h.keymap[action]
	if !ok {
		return false
	}
	for _, k := range keys {
		if h.keyIsPressed(k) {
			return true
		}
	}
	return false
}

func (h *Handler) keyIsJustPressed(k Key) bool {
	switch k.kind {
	case keyTouch:
		if k.code == int(touchTap) {
			return h.sys.touchHasTap
		}
		return false
	case keyGamepad:
		return h.gamepadKeyIsJustPressed(k)
	case keyGamepadLeftStick:
		return h.gamepadStickIsJustPressed(stickCode(k.code), ebiten.StandardGamepadAxisLeftStickHorizontal, ebiten.StandardGamepadAxisLeftStickVertical)
	case keyGamepadRightStick:
		return h.gamepadStickIsJustPressed(stickCode(k.code), ebiten.StandardGamepadAxisRightStickHorizontal, ebiten.StandardGamepadAxisRightStickVertical)
	case keyMouse:
		return inpututil.IsMouseButtonJustPressed(ebiten.MouseButton(k.code))
	case keyKeyboardWithCtrl:
		return ebiten.IsKeyPressed(ebiten.KeyControl) &&
			inpututil.IsKeyJustPressed(ebiten.Key(k.code))
	case keyKeyboardWithShift:
		return ebiten.IsKeyPressed(ebiten.KeyShift) &&
			inpututil.IsKeyJustPressed(ebiten.Key(k.code))
	default:
		return inpututil.IsKeyJustPressed(ebiten.Key(k.code))
	}
}

func (h *Handler) keyIsPressed(k Key) bool {
	switch k.kind {
	case keyGamepad:
		return h.gamepadKeyIsPressed(k)
	case keyGamepadLeftStick:
		return h.gamepadStickIsPressed(stickCode(k.code), ebiten.StandardGamepadAxisLeftStickHorizontal, ebiten.StandardGamepadAxisLeftStickVertical)
	case keyGamepadRightStick:
		return h.gamepadStickIsPressed(stickCode(k.code), ebiten.StandardGamepadAxisRightStickHorizontal, ebiten.StandardGamepadAxisRightStickVertical)
	case keyMouse:
		return ebiten.IsMouseButtonPressed(ebiten.MouseButton(k.code))
	case keyKeyboardWithCtrl:
		return ebiten.IsKeyPressed(ebiten.KeyControl) &&
			ebiten.IsKeyPressed(ebiten.Key(k.code))
	case keyKeyboardWithShift:
		return ebiten.IsKeyPressed(ebiten.KeyShift) &&
			ebiten.IsKeyPressed(ebiten.Key(k.code))
	default:
		return ebiten.IsKeyPressed(ebiten.Key(k.code))
	}
}

func (h *Handler) isDPadAxisActive(code int, vec Vec) bool {
	switch ebiten.StandardGamepadButton(code) {
	case ebiten.StandardGamepadButtonLeftTop:
		return vec.Y == -1
	case ebiten.StandardGamepadButtonLeftRight:
		return vec.X == 1
	case ebiten.StandardGamepadButtonLeftBottom:
		return vec.Y == 1
	case ebiten.StandardGamepadButtonLeftLeft:
		return vec.X == -1
	}
	return false
}

func (h *Handler) gamepadKeyIsJustPressed(k Key) bool {
	if h.gamepadInfo().model == gamepadStandard {
		return inpututil.IsStandardGamepadButtonJustPressed(ebiten.GamepadID(h.id), ebiten.StandardGamepadButton(k.code))
	}
	if h.gamepadInfo().model == gamepadFirefoxXinput && isDPadButton(k.code) {
		return !h.isDPadAxisActive(k.code, h.getStickPrevVec(6, 7)) &&
			h.isDPadAxisActive(k.code, h.getStickVec(6, 7))
	}
	return inpututil.IsGamepadButtonJustPressed(ebiten.GamepadID(h.id), h.mappedGamepadKey(k.code))
}

func (h *Handler) gamepadKeyIsPressed(k Key) bool {
	if h.gamepadInfo().model == gamepadStandard {
		return ebiten.IsStandardGamepadButtonPressed(ebiten.GamepadID(h.id), ebiten.StandardGamepadButton(k.code))
	}
	if h.gamepadInfo().model == gamepadFirefoxXinput && isDPadButton(k.code) {
		return h.isDPadAxisActive(k.code, h.getStickVec(6, 7))
	}
	return ebiten.IsGamepadButtonPressed(ebiten.GamepadID(h.id), h.mappedGamepadKey(k.code))
}

func (h *Handler) gamepadStickIsActive(code stickCode, vec Vec) bool {
	if vecLen(vec) < 0.5 {
		return false
	}
	switch code {
	case stickUp:
		angle := angleNormalized(vecAngle(vec))
		return angle > (math.Pi+math.Pi/4) && angle <= (2*math.Pi-math.Pi/4)
	case stickRight:
		angle := angleNormalized(vecAngle(vec))
		return angle <= (math.Pi/4) || angle > (2*math.Pi-math.Pi/4)
	case stickDown:
		angle := angleNormalized(vecAngle(vec))
		return angle > (math.Pi/4) && angle <= (math.Pi-math.Pi/4)
	case stickLeft:
		angle := angleNormalized(vecAngle(vec))
		return angle > (math.Pi-math.Pi/4) && angle <= (math.Pi+math.Pi/4)
	}
	return false
}

func (h *Handler) gamepadStickIsJustPressed(code stickCode, axis1, axis2 ebiten.StandardGamepadAxis) bool {
	return !h.gamepadStickIsActive(code, h.getStickPrevVec(axis1, axis2)) &&
		h.gamepadStickIsActive(code, h.getStickVec(int(axis1), int(axis2)))
}

func (h *Handler) gamepadStickIsPressed(code stickCode, axis1, axis2 ebiten.StandardGamepadAxis) bool {
	vec := h.getStickVec(int(axis1), int(axis2))
	return h.gamepadStickIsActive(code, vec)
}

func (h *Handler) getStickPrevVec(axis1, axis2 ebiten.StandardGamepadAxis) Vec {
	return Vec{
		X: h.gamepadInfo().prevAxisValues[axis1],
		Y: h.gamepadInfo().prevAxisValues[axis2],
	}
}

func (h *Handler) getStickVec(axis1, axis2 int) Vec {
	return Vec{
		X: h.gamepadInfo().axisValues[axis1],
		Y: h.gamepadInfo().axisValues[axis2],
	}
}

func (h *Handler) gamepadInfo() *gamepadInfo {
	return &h.sys.gamepadInfo[h.id]
}

func (h *Handler) mappedGamepadKey(keyCode int) ebiten.GamepadButton {
	b := ebiten.StandardGamepadButton(keyCode)
	switch h.gamepadInfo().model {
	case gamepadMicront:
		return microntToXbox(b)
	case gamepadFirefoxXinput:
		return firefoxXinputToXbox(b)
	default:
		return ebiten.GamepadButton(keyCode)
	}
}
