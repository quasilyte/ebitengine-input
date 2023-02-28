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
	id     uint8
	keymap Keymap
	sys    *System
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
func (h *Handler) DefaultInputMask() DeviceKind {
	if h.GamepadConnected() {
		return GamepadDevice
	}
	return KeyboardDevice | MouseDevice
}

// EmitKeyEvent sends given key event into the input system.
//
// The event is emitted from the perspective of this handler,
// so the gamepad events will be handled properly in the multi-device context.
//
// Note: simulated events are only visible after the next System.Update() call.
//
// See SimulatedKeyEvent documentation for more info.
//
// Experimental: this is a part of virtual input API, which is not stable yet.
func (h *Handler) EmitKeyEvent(e SimulatedKeyEvent) {
	h.sys.pendingEvents = append(h.sys.pendingEvents, simulatedEvent{
		code:     e.Key.code,
		keyKind:  e.Key.kind,
		playerID: h.id,
		pos:      e.Pos,
	})
}

// EmitEvent activates the given action for the player.
// Only the handlers with the same player ID will discover this action.
//
// Note: simulated events are only visible after the next System.Update() call.p
//
// See SimulatedAction documentation for more info.
//
// Experimental: this is a part of virtual input API, which is not stable yet.
func (h *Handler) EmitEvent(e SimulatedAction) {
	h.sys.pendingEvents = append(h.sys.pendingEvents, simulatedEvent{
		code:     int(e.Action),
		keyKind:  keySimulated,
		playerID: h.id,
		pos:      e.Pos,
		startPos: e.StartPos,
	})
}

// ActionKeyNames returns a list of key names associated by this action.
//
// It filters the results by a given input device mask.
// If you want to include all input device keys, use AnyDevice value.
//
// This function is useful when you want to display a list of keys
// the player should press in order to activate some action.
//
// The filtering is useful to avoid listing the unrelated options.
// For example, if player uses the gamepad, it could be weird to
// show keyboard options listed. For the simple cases, you can use
// DefaultInputMask() method to get the mask that will try to avoid
// that situation. See its comment to learn more.
//
// Keys with modifiers will have them listed too.
// Modifiers are separated by "+".
// A "k" keyboard key with ctrl modifier will have a "ctrl+k" name.
func (h *Handler) ActionKeyNames(action Action, mask DeviceKind) []string {
	keys, ok := h.keymap[action]
	if !ok {
		return nil
	}
	gamepadConnected := h.GamepadConnected()
	result := make([]string, 0, len(keys))
	for _, k := range keys {
		enabled := true
		ctrlMod := false
		shiftMod := false
		switch k.kind {
		case keyKeyboardWithCtrlShift:
			ctrlMod = true
			shiftMod = true
			enabled = mask&KeyboardDevice != 0
		case keyKeyboardWithCtrl:
			ctrlMod = true
			enabled = mask&KeyboardDevice != 0
		case keyKeyboardWithShift:
			shiftMod = true
			enabled = mask&KeyboardDevice != 0
		case keyKeyboard:
			enabled = mask&KeyboardDevice != 0
		case keyMouseWithCtrlShift:
			ctrlMod = true
			shiftMod = true
			enabled = mask&MouseDevice != 0
		case keyMouseWithCtrl:
			ctrlMod = true
			enabled = mask&MouseDevice != 0
		case keyMouseWithShift:
			shiftMod = true
			enabled = mask&MouseDevice != 0
		case keyMouse:
			enabled = mask&MouseDevice != 0
		case keyGamepad, keyGamepadLeftStick, keyGamepadRightStick, keyGamepadStickMotion:
			enabled = gamepadConnected && (mask&GamepadDevice != 0)
		case keyTouch, keyTouchDrag:
			enabled = h.sys.touchEnabled && (mask&TouchDevice != 0)
		}
		if enabled {
			name := k.name
			if shiftMod {
				name = "shift+" + name
			}
			if ctrlMod {
				name = "ctrl+" + name
			}
			result = append(result, name)
		}
	}
	return result
}

// JustPressedActionInfo is like ActionIsJustPressed, but with more information.
//
// The first return value will hold the extra event info.
// The second return value is false if given action is not activated.
//
// See EventInfo comment to learn more.
func (h *Handler) JustPressedActionInfo(action Action) (EventInfo, bool) {
	keys, ok := h.keymap[action]
	if !ok {
		return EventInfo{}, false
	}
	for _, k := range keys {
		if info, status := h.pressedSimulatedKeyInfo(true, k); status == bool3true {
			return info, true
		}
		if !h.keyIsJustPressed(k) {
			continue
		}
		var info EventInfo
		info.kind = k.kind
		info.hasPos = keyHasPos(k.kind)
		info.Pos = h.getKeyPos(k)
		info.StartPos = h.getKeyStartPos(k)
		return info, true
	}
	if h.sys.hasSimulatedActions {
		info, status := h.pressedSimulatedKeyInfo(true, Key{
			code: int(action),
			kind: keySimulated,
		})
		return info, status == bool3true
	}
	return EventInfo{}, false
}

// PressedActionInfo is like ActionIsPressed, but with more information.
//
// The first return value will hold the extra event info.
// The second return value is false if given action is not activated.
//
// See EventInfo comment to learn more.
func (h *Handler) PressedActionInfo(action Action) (EventInfo, bool) {
	keys, ok := h.keymap[action]
	if !ok {
		return EventInfo{}, false
	}
	for _, k := range keys {
		if info, status := h.pressedSimulatedKeyInfo(false, k); status == bool3true {
			return info, true
		}
		if !h.keyIsPressed(k) {
			continue
		}
		var info EventInfo
		info.kind = k.kind
		info.hasPos = keyHasPos(k.kind)
		info.Pos = h.getKeyPos(k)
		info.StartPos = h.getKeyStartPos(k)
		return info, true
	}
	return EventInfo{}, false
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
		if len(h.sys.simulatedEvents) != 0 {
			// We want to avoid a situation when simulated input
			// things that the key is still being pressed and then
			// receive a real input from the bottom of ebitenutil that
			// this key was actually "just pressed". To avoid that,
			// we skip checking the real input if simulated input still
			// holds that button down. This is why we need a bool3 here.
			_, isPressed := h.pressedSimulatedKeyInfo(true, k)
			if isPressed != bool3unset {
				return isPressed == bool3true
			}
		}
		if h.keyIsJustPressed(k) {
			return true
		}
	}
	if h.sys.hasSimulatedActions {
		_, isPressed := h.pressedSimulatedKeyInfo(true, Key{
			code: int(action),
			kind: keySimulated,
		})
		return isPressed == bool3true
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
		if len(h.sys.simulatedEvents) != 0 && h.simulatedKeyIsPressed(k) {
			return true
		}
		if h.keyIsPressed(k) {
			return true
		}
	}
	if h.sys.hasSimulatedActions {
		return h.simulatedKeyIsPressed(Key{
			code: int(action),
			kind: keySimulated,
		})
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
	case keyTouchDrag:
		return h.sys.touchJustHadDrag
	case keyGamepad:
		return h.gamepadKeyIsJustPressed(k)
	case keyGamepadLeftStick:
		return h.gamepadStickIsJustPressed(stickCode(k.code), ebiten.StandardGamepadAxisLeftStickHorizontal, ebiten.StandardGamepadAxisLeftStickVertical)
	case keyGamepadRightStick:
		return h.gamepadStickIsJustPressed(stickCode(k.code), ebiten.StandardGamepadAxisRightStickHorizontal, ebiten.StandardGamepadAxisRightStickVertical)
	case keyGamepadStickMotion:
		return h.gamepadStickMotionIsJustPressed(stickCode(k.code))
	case keyMouse:
		return inpututil.IsMouseButtonJustPressed(ebiten.MouseButton(k.code))
	case keyMouseWithCtrl:
		return ebiten.IsKeyPressed(ebiten.KeyControl) &&
			inpututil.IsMouseButtonJustPressed(ebiten.MouseButton(k.code))
	case keyMouseWithShift:
		return ebiten.IsKeyPressed(ebiten.KeyShift) &&
			inpututil.IsMouseButtonJustPressed(ebiten.MouseButton(k.code))
	case keyMouseWithCtrlShift:
		return ebiten.IsKeyPressed(ebiten.KeyControl) &&
			ebiten.IsKeyPressed(ebiten.KeyShift) &&
			inpututil.IsMouseButtonJustPressed(ebiten.MouseButton(k.code))
	case keyKeyboardWithCtrl:
		return ebiten.IsKeyPressed(ebiten.KeyControl) &&
			inpututil.IsKeyJustPressed(ebiten.Key(k.code))
	case keyKeyboardWithShift:
		return ebiten.IsKeyPressed(ebiten.KeyShift) &&
			inpututil.IsKeyJustPressed(ebiten.Key(k.code))
	case keyKeyboardWithCtrlShift:
		return ebiten.IsKeyPressed(ebiten.KeyControl) &&
			ebiten.IsKeyPressed(ebiten.KeyShift) &&
			inpututil.IsKeyJustPressed(ebiten.Key(k.code))
	case keyWheel:
		return h.wheelIsJustPressed(wheelCode(k.code))
	default:
		return inpututil.IsKeyJustPressed(ebiten.Key(k.code))
	}
}

func (h *Handler) getKeyStartPos(k Key) Vec {
	var result Vec
	switch k.kind {
	case keyTouchDrag:
		result = h.sys.touchStartPos
	}
	return result
}

func (h *Handler) getKeyPos(k Key) Vec {
	var result Vec
	switch k.kind {
	case keyMouse, keyMouseWithCtrl, keyMouseWithShift, keyMouseWithCtrlShift:
		result = h.sys.cursorPos
	case keyTouch:
		result = h.sys.touchTapPos
	case keyTouchDrag:
		result = h.sys.touchDragPos
	case keyWheel:
		result = h.sys.wheel
	case keyGamepadStickMotion:
		axis1, axis2 := h.getStickAxes(stickCode(k.code))
		result = h.getStickVec(axis1, axis2)
	}
	return result
}

func (h *Handler) keyIsPressed(k Key) bool {
	switch k.kind {
	case keyTouch:
		if k.code == int(touchTap) {
			return h.sys.touchHasTap
		}
		return false
	case keyTouchDrag:
		return h.sys.touchHasDrag
	case keyGamepad:
		return h.gamepadKeyIsPressed(k)
	case keyGamepadLeftStick:
		return h.gamepadStickIsPressed(stickCode(k.code), ebiten.StandardGamepadAxisLeftStickHorizontal, ebiten.StandardGamepadAxisLeftStickVertical)
	case keyGamepadRightStick:
		return h.gamepadStickIsPressed(stickCode(k.code), ebiten.StandardGamepadAxisRightStickHorizontal, ebiten.StandardGamepadAxisRightStickVertical)
	case keyGamepadStickMotion:
		return h.gamepadStickMotionIsPressed(stickCode(k.code))
	case keyMouse:
		return ebiten.IsMouseButtonPressed(ebiten.MouseButton(k.code))
	case keyMouseWithCtrl:
		return ebiten.IsKeyPressed(ebiten.KeyControl) &&
			ebiten.IsMouseButtonPressed(ebiten.MouseButton(k.code))
	case keyMouseWithShift:
		return ebiten.IsKeyPressed(ebiten.KeyShift) &&
			ebiten.IsMouseButtonPressed(ebiten.MouseButton(k.code))
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

func (h *Handler) eventSliceFind(slice []simulatedEvent, k Key) int {
	for i, e := range slice {
		if e.code == k.code && e.keyKind == k.kind {
			if keyNeedID(e.keyKind) && e.playerID != h.id {
				continue
			}
			return i
		}
	}
	return -1
}

func (h *Handler) eventSliceContains(slice []simulatedEvent, k Key) bool {
	return h.eventSliceFind(slice, k) != -1
}

func (h *Handler) pressedSimulatedKeyInfo(justPressed bool, k Key) (EventInfo, bool3) {
	var info EventInfo
	i := h.eventSliceFind(h.sys.simulatedEvents, k)
	if i != -1 {
		if justPressed && h.eventSliceContains(h.sys.prevSimulatedEvents, k) {
			return info, bool3false
		}
		info.Pos = h.sys.simulatedEvents[i].pos
		info.StartPos = h.sys.simulatedEvents[i].startPos
		info.kind = k.kind
		info.hasPos = keyHasPos(k.kind)
		return info, bool3true
	}
	return info, bool3unset
}

func (h *Handler) simulatedKeyIsPressed(k Key) bool {
	return h.eventSliceContains(h.sys.simulatedEvents, k)
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

func (h *Handler) wheelIsJustPressed(code wheelCode) bool {
	switch code {
	case wheelDown:
		return h.sys.wheel.Y > 0
	case wheelUp:
		return h.sys.wheel.Y < 0
	case wheelVertical:
		return h.sys.wheel.Y != 0
	default:
		return false
	}
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
	// Allow some axis overlap to emulate double direction pressing,
	// like with D-pad buttons.
	const overlap float64 = math.Pi / 7
	switch code {
	case stickUp:
		angle := angleNormalized(vecAngle(vec))
		return angle > (math.Pi+math.Pi/4)-overlap && angle <= (2*math.Pi-math.Pi/4)+overlap
	case stickRight:
		angle := angleNormalized(vecAngle(vec))
		return angle <= (math.Pi/4)+overlap || angle > (2*math.Pi-math.Pi/4)-overlap
	case stickDown:
		angle := angleNormalized(vecAngle(vec))
		return angle > (math.Pi/4)-overlap && angle <= (math.Pi-math.Pi/4)+overlap
	case stickLeft:
		angle := angleNormalized(vecAngle(vec))
		return angle > (math.Pi-math.Pi/4)-overlap && angle <= (math.Pi+math.Pi/4)+overlap
	}
	return false
}

func (h *Handler) gamepadStickIsJustPressed(code stickCode, axis1, axis2 ebiten.StandardGamepadAxis) bool {
	return !h.gamepadStickIsActive(code, h.getStickPrevVec(int(axis1), int(axis2))) &&
		h.gamepadStickIsActive(code, h.getStickVec(int(axis1), int(axis2)))
}

func (h *Handler) getStickAxes(code stickCode) (int, int) {
	var axis1 int
	var axis2 int
	if code == stickLeft {
		axis1 = int(ebiten.StandardGamepadAxisLeftStickHorizontal)
		axis2 = int(ebiten.StandardGamepadAxisLeftStickVertical)
	} else {
		axis1 = int(ebiten.StandardGamepadAxisRightStickHorizontal)
		axis2 = int(ebiten.StandardGamepadAxisRightStickVertical)
	}
	return axis1, axis2
}

func (h *Handler) gamepadStickMotionIsJustPressed(code stickCode) bool {
	return !h.gamepadStickMotionIsActive(h.getStickPrevVec(h.getStickAxes(code))) &&
		h.gamepadStickMotionIsActive(h.getStickVec(h.getStickAxes(code)))
}

func (h *Handler) gamepadStickMotionIsPressed(code stickCode) bool {
	return h.gamepadStickMotionIsActive(h.getStickVec(h.getStickAxes(code)))
}

func (h *Handler) gamepadStickMotionIsActive(vec Vec) bool {
	// Some gamepads could register a slight movement all the time,
	// even if the stick is in its home position.
	// This min sensitivity should probably be configurable.
	// My gamepads may have false positive activations with a
	// value lower than 0.03; we're using 0.055 here just to be safe.
	return math.Abs(vec.X)+math.Abs(vec.Y) >= 0.055
}

func (h *Handler) gamepadStickIsPressed(code stickCode, axis1, axis2 ebiten.StandardGamepadAxis) bool {
	vec := h.getStickVec(int(axis1), int(axis2))
	return h.gamepadStickIsActive(code, vec)
}

func (h *Handler) getStickPrevVec(axis1, axis2 int) Vec {
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
