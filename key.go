package input

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Key represents an input method that can be used to activate Action.
// Key could be a keyboard key, a gamepad key, a mouse button, etc.
//
// Use the predefined global vars like KeyMouseLeft and KeyTab to create a Keymap.
type Key struct {
	code int
	kind keyKind
	name string
}

type KeyModifier uint8

const (
	ModControl KeyModifier = iota
	ModShift
	ModControlShift
)

// KeyWithModifier turns k into a combined modifier+k key.
// For instance, KeyUp+ModControl will trigger an action
// only if both of these keys are being pressed.
func KeyWithModifier(k Key, mod KeyModifier) Key {
	switch k.kind {
	case keyKeyboard:
		switch mod {
		case ModControl:
			k.kind = keyKeyboardWithCtrl
		case ModShift:
			k.kind = keyKeyboardWithShift
		case ModControlShift:
			k.kind = keyKeyboardWithCtrlShift
		default:
			panic("unexpected key modifier")
		}
	case keyMouse:
		switch mod {
		case ModControl:
			k.kind = keyMouseWithCtrl
		case ModShift:
			k.kind = keyMouseWithShift
		case ModControlShift:
			k.kind = keyMouseWithCtrlShift
		default:
			panic("unexpected key modifier")
		}
	default:
		panic("only keyboard and mouse keys support modifiers")
	}
	return k
}

// Wheel keys.
var (
	KeyWheelUp       = Key{code: int(wheelUp), kind: keyWheel, name: "wheel_up"}
	KeyWheelDown     = Key{code: int(wheelDown), kind: keyWheel, name: "wheel_down"}
	KeyWheelVertical = Key{code: int(wheelVertical), kind: keyWheel, name: "wheel_vertical"}
)

// Mouse keys.
var (
	KeyMouseLeft   = Key{code: int(ebiten.MouseButtonLeft), kind: keyMouse, name: "mouse_left_button"}
	KeyMouseRight  = Key{code: int(ebiten.MouseButtonRight), kind: keyMouse, name: "mouse_right_button"}
	KeyMouseMiddle = Key{code: int(ebiten.MouseButtonMiddle), kind: keyMouse, name: "mouse_middle_button"}
)

// Touch keys.
var (
	KeyTouchTap = Key{code: int(touchTap), kind: keyTouch, name: "screen_tap"}
)

// Keyboard keys.
var (
	KeyLeft  = Key{code: int(ebiten.KeyLeft), name: "left"}
	KeyRight = Key{code: int(ebiten.KeyRight), name: "right"}
	KeyUp    = Key{code: int(ebiten.KeyUp), name: "up"}
	KeyDown  = Key{code: int(ebiten.KeyDown), name: "down"}

	Key0     = Key{code: int(ebiten.Key0), name: "0"}
	Key1     = Key{code: int(ebiten.Key1), name: "1"}
	Key2     = Key{code: int(ebiten.Key2), name: "2"}
	Key3     = Key{code: int(ebiten.Key3), name: "3"}
	Key4     = Key{code: int(ebiten.Key4), name: "4"}
	Key5     = Key{code: int(ebiten.Key5), name: "5"}
	Key6     = Key{code: int(ebiten.Key6), name: "6"}
	Key7     = Key{code: int(ebiten.Key7), name: "7"}
	Key8     = Key{code: int(ebiten.Key8), name: "8"}
	Key9     = Key{code: int(ebiten.Key9), name: "9"}
	KeyMinus = Key{code: int(ebiten.KeyMinus), name: "minus"}
	KeyEqual = Key{code: int(ebiten.KeyEqual), name: "equal"}

	KeyBackquote = Key{code: int(ebiten.KeyBackquote), name: "backquote"}

	KeyQ = Key{code: int(ebiten.KeyQ), name: "q"}
	KeyW = Key{code: int(ebiten.KeyW), name: "w"}
	KeyE = Key{code: int(ebiten.KeyE), name: "e"}
	KeyR = Key{code: int(ebiten.KeyR), name: "r"}
	KeyT = Key{code: int(ebiten.KeyT), name: "t"}
	KeyY = Key{code: int(ebiten.KeyY), name: "y"}
	KeyU = Key{code: int(ebiten.KeyY), name: "u"}
	KeyI = Key{code: int(ebiten.KeyY), name: "i"}
	KeyO = Key{code: int(ebiten.KeyY), name: "o"}
	KeyP = Key{code: int(ebiten.KeyY), name: "p"}
	KeyA = Key{code: int(ebiten.KeyA), name: "a"}
	KeyS = Key{code: int(ebiten.KeyS), name: "s"}
	KeyD = Key{code: int(ebiten.KeyD), name: "d"}
	KeyF = Key{code: int(ebiten.KeyF), name: "f"}
	KeyG = Key{code: int(ebiten.KeyG), name: "g"}
	KeyH = Key{code: int(ebiten.KeyH), name: "h"}
	KeyJ = Key{code: int(ebiten.KeyJ), name: "j"}
	KeyK = Key{code: int(ebiten.KeyK), name: "k"}
	KeyL = Key{code: int(ebiten.KeyL), name: "l"}
	KeyZ = Key{code: int(ebiten.KeyZ), name: "z"}
	KeyX = Key{code: int(ebiten.KeyX), name: "x"}
	KeyC = Key{code: int(ebiten.KeyC), name: "c"}
	KeyV = Key{code: int(ebiten.KeyV), name: "v"}
	KeyB = Key{code: int(ebiten.KeyB), name: "b"}
	KeyN = Key{code: int(ebiten.KeyN), name: "n"}
	KeyM = Key{code: int(ebiten.KeyM), name: "m"}

	KeyBackspace    = Key{code: int(ebiten.KeyBackspace), name: "backspace"}
	KeyComma        = Key{code: int(ebiten.KeyComma), name: "comma"}
	KeyControl      = Key{code: int(ebiten.KeyControl), name: "control"}
	KeyControlLeft  = Key{code: int(ebiten.KeyControlLeft), name: "control_left"}
	KeyControlRight = Key{code: int(ebiten.KeyControlRight), name: "control_right"}
	KeyDelete       = Key{code: int(ebiten.KeyDelete), name: "delete"}
	KeyEnd          = Key{code: int(ebiten.KeyEnd), name: "end"}
	KeyEnter        = Key{code: int(ebiten.KeyEnter), name: "enter"}
	KeyEscape       = Key{code: int(ebiten.KeyEscape), name: "escape"}
	KeyHome         = Key{code: int(ebiten.KeyHome), name: "home"}
	KeyInsert       = Key{code: int(ebiten.KeyInsert), name: "insert"}
	KeyPageDown     = Key{code: int(ebiten.KeyPageDown), name: "page_down"}
	KeyPageUp       = Key{code: int(ebiten.KeyPageUp), name: "page_up"}
	KeyPeriod       = Key{code: int(ebiten.KeyPeriod), name: "period"}
	KeyShift        = Key{code: int(ebiten.KeyShift), name: "shift"}
	KeyShiftLeft    = Key{code: int(ebiten.KeyShiftLeft), name: "shift_left"}
	KeyShiftRight   = Key{code: int(ebiten.KeyShiftRight), name: "shift_right"}
	KeySlash        = Key{code: int(ebiten.KeySlash), name: "slash"}
	KeySpace        = Key{code: int(ebiten.KeySpace), name: "space"}
	KeyTab          = Key{code: int(ebiten.KeyTab), name: "tab"}

	KeyNum0 = Key{code: int(ebiten.KeyNumpad0), name: "numpad_0"}
	KeyNum1 = Key{code: int(ebiten.KeyNumpad1), name: "numpad_1"}
	KeyNum2 = Key{code: int(ebiten.KeyNumpad2), name: "numpad_2"}
	KeyNum3 = Key{code: int(ebiten.KeyNumpad3), name: "numpad_3"}
	KeyNum4 = Key{code: int(ebiten.KeyNumpad4), name: "numpad_4"}
	KeyNum5 = Key{code: int(ebiten.KeyNumpad5), name: "numpad_5"}
	KeyNum6 = Key{code: int(ebiten.KeyNumpad6), name: "numpad_6"}
	KeyNum7 = Key{code: int(ebiten.KeyNumpad7), name: "numpad_7"}
	KeyNum8 = Key{code: int(ebiten.KeyNumpad8), name: "numpad_8"}
	KeyNum9 = Key{code: int(ebiten.KeyNumpad9), name: "numpad_9"}

	KeyNumAdd      = Key{code: int(ebiten.KeyNumpadAdd), name: "numpad_add"}
	KeyNumDivide   = Key{code: int(ebiten.KeyNumpadDivide), name: "numpad_divide"}
	KeyNumEnter    = Key{code: int(ebiten.KeyNumpadEnter), name: "numpad_enter"}
	KeyNumMultiply = Key{code: int(ebiten.KeyNumpadMultiply), name: "numpad_multiply"}
	KeyNumPeriod   = Key{code: int(ebiten.KeyNumpadDecimal), name: "numpad_period"}
	KeyNumSubtract = Key{code: int(ebiten.KeyNumpadSubtract), name: "numpad_subtract"}
)

// Gamepad keys.
var (
	KeyGamepadStart  = Key{code: int(ebiten.StandardGamepadButtonCenterRight), kind: keyGamepad, name: "gamepad_start"}
	KeyGamepadSelect = Key{code: int(ebiten.StandardGamepadButtonCenterLeft), kind: keyGamepad, name: "gamepad_select"}
	KeyGamepadMiddle = Key{code: int(ebiten.StandardGamepadButtonCenterCenter), kind: keyGamepad, name: "gamepad_middle"}

	KeyGamepadUp    = Key{code: int(ebiten.StandardGamepadButtonLeftTop), kind: keyGamepad, name: "gamepad_up"}
	KeyGamepadRight = Key{code: int(ebiten.StandardGamepadButtonLeftRight), kind: keyGamepad, name: "gamepad_right"}
	KeyGamepadDown  = Key{code: int(ebiten.StandardGamepadButtonLeftBottom), kind: keyGamepad, name: "gamepad_down"}
	KeyGamepadLeft  = Key{code: int(ebiten.StandardGamepadButtonLeftLeft), kind: keyGamepad, name: "gamepad_left"}

	KeyGamepadLStickUp    = Key{code: int(stickUp), kind: keyGamepadLeftStick, name: "gamepad_lstick_up"}
	KeyGamepadLStickRight = Key{code: int(stickRight), kind: keyGamepadLeftStick, name: "gamepad_lstick_right"}
	KeyGamepadLStickDown  = Key{code: int(stickDown), kind: keyGamepadLeftStick, name: "gamepad_lstick_down"}
	KeyGamepadLStickLeft  = Key{code: int(stickLeft), kind: keyGamepadLeftStick, name: "gamepad_lstick_left"}
	KeyGamepadRStickUp    = Key{code: int(stickUp), kind: keyGamepadRightStick, name: "gamepad_rstick_up"}
	KeyGamepadRStickRight = Key{code: int(stickRight), kind: keyGamepadRightStick, name: "gamepad_rstick_right"}
	KeyGamepadRStickDown  = Key{code: int(stickDown), kind: keyGamepadRightStick, name: "gamepad_rstick_down"}
	KeyGamepadRStickLeft  = Key{code: int(stickLeft), kind: keyGamepadRightStick, name: "gamepad_rstick_left"}

	KeyGamepadA = Key{code: int(ebiten.StandardGamepadButtonRightBottom), kind: keyGamepad, name: "gamepad_a"}
	KeyGamepadB = Key{code: int(ebiten.StandardGamepadButtonRightRight), kind: keyGamepad, name: "gamepad_b"}
	KeyGamepadX = Key{code: int(ebiten.StandardGamepadButtonRightLeft), kind: keyGamepad, name: "gamepad_x"}
	KeyGamepadY = Key{code: int(ebiten.StandardGamepadButtonRightTop), kind: keyGamepad, name: "gamepad_y"}

	KeyGamepadL1 = Key{code: int(ebiten.StandardGamepadButtonFrontTopLeft), kind: keyGamepad, name: "gamepad_l1"}
	KeyGamepadL2 = Key{code: int(ebiten.StandardGamepadButtonFrontBottomLeft), kind: keyGamepad, name: "gamepad_l2"}
	KeyGamepadR1 = Key{code: int(ebiten.StandardGamepadButtonFrontTopRight), kind: keyGamepad, name: "gamepad_r1"}
	KeyGamepadR2 = Key{code: int(ebiten.StandardGamepadButtonFrontBottomRight), kind: keyGamepad, name: "gamepad_r2"}
)
