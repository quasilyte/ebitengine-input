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
	KeyLeft  Key = Key{code: int(ebiten.KeyLeft), name: "left"}
	KeyRight Key = Key{code: int(ebiten.KeyRight), name: "right"}
	KeyUp    Key = Key{code: int(ebiten.KeyUp), name: "up"}
	KeyDown  Key = Key{code: int(ebiten.KeyDown), name: "down"}

	Key0 Key = Key{code: int(ebiten.Key0), name: "0"}
	Key1 Key = Key{code: int(ebiten.Key1), name: "1"}
	Key2 Key = Key{code: int(ebiten.Key2), name: "2"}
	Key3 Key = Key{code: int(ebiten.Key3), name: "3"}
	Key4 Key = Key{code: int(ebiten.Key4), name: "4"}
	Key5 Key = Key{code: int(ebiten.Key5), name: "5"}
	Key6 Key = Key{code: int(ebiten.Key6), name: "6"}
	Key7 Key = Key{code: int(ebiten.Key7), name: "7"}
	Key8 Key = Key{code: int(ebiten.Key8), name: "8"}
	Key9 Key = Key{code: int(ebiten.Key9), name: "9"}

	KeyQ Key = Key{code: int(ebiten.KeyQ), name: "q"}
	KeyW Key = Key{code: int(ebiten.KeyW), name: "w"}
	KeyE Key = Key{code: int(ebiten.KeyE), name: "e"}
	KeyR Key = Key{code: int(ebiten.KeyR), name: "r"}
	KeyT Key = Key{code: int(ebiten.KeyT), name: "t"}
	KeyY Key = Key{code: int(ebiten.KeyY), name: "y"}
	KeyU Key = Key{code: int(ebiten.KeyY), name: "u"}
	KeyI Key = Key{code: int(ebiten.KeyY), name: "i"}
	KeyO Key = Key{code: int(ebiten.KeyY), name: "o"}
	KeyP Key = Key{code: int(ebiten.KeyY), name: "p"}
	KeyA Key = Key{code: int(ebiten.KeyA), name: "a"}
	KeyS Key = Key{code: int(ebiten.KeyS), name: "s"}
	KeyD Key = Key{code: int(ebiten.KeyD), name: "d"}
	KeyF Key = Key{code: int(ebiten.KeyD), name: "f"}
	KeyG Key = Key{code: int(ebiten.KeyD), name: "g"}
	KeyH Key = Key{code: int(ebiten.KeyD), name: "h"}
	KeyJ Key = Key{code: int(ebiten.KeyD), name: "j"}
	KeyK Key = Key{code: int(ebiten.KeyD), name: "k"}
	KeyL Key = Key{code: int(ebiten.KeyD), name: "l"}
	KeyZ Key = Key{code: int(ebiten.KeyD), name: "z"}
	KeyX Key = Key{code: int(ebiten.KeyD), name: "x"}
	KeyC Key = Key{code: int(ebiten.KeyD), name: "c"}
	KeyV Key = Key{code: int(ebiten.KeyD), name: "v"}
	KeyB Key = Key{code: int(ebiten.KeyD), name: "b"}
	KeyN Key = Key{code: int(ebiten.KeyD), name: "n"}
	KeyM Key = Key{code: int(ebiten.KeyD), name: "m"}

	KeyEscape Key = Key{code: int(ebiten.KeyEscape), name: "escape"}
	KeyEnter  Key = Key{code: int(ebiten.KeyEnter), name: "enter"}
	KeyTab    Key = Key{code: int(ebiten.KeyTab), name: "tab"}
	KeySpace  Key = Key{code: int(ebiten.KeySpace), name: "space"}
)

// Gamepad keys.
var (
	KeyGamepadStart  Key = Key{code: int(ebiten.StandardGamepadButtonCenterRight), kind: keyGamepad, name: "gamepad_start"}
	KeyGamepadSelect Key = Key{code: int(ebiten.StandardGamepadButtonCenterLeft), kind: keyGamepad, name: "gamepad_select"}
	KeyGamepadMiddle Key = Key{code: int(ebiten.StandardGamepadButtonCenterCenter), kind: keyGamepad, name: "gamepad_middle"}

	KeyGamepadUp    Key = Key{code: int(ebiten.StandardGamepadButtonLeftTop), kind: keyGamepad, name: "gamepad_up"}
	KeyGamepadRight Key = Key{code: int(ebiten.StandardGamepadButtonLeftRight), kind: keyGamepad, name: "gamepad_right"}
	KeyGamepadDown  Key = Key{code: int(ebiten.StandardGamepadButtonLeftBottom), kind: keyGamepad, name: "gamepad_down"}
	KeyGamepadLeft  Key = Key{code: int(ebiten.StandardGamepadButtonLeftLeft), kind: keyGamepad, name: "gamepad_left"}

	KeyGamepadLStickUp    = Key{code: int(stickUp), kind: keyGamepadLeftStick, name: "gamepad_lstick_up"}
	KeyGamepadLStickRight = Key{code: int(stickRight), kind: keyGamepadLeftStick, name: "gamepad_lstick_right"}
	KeyGamepadLStickDown  = Key{code: int(stickDown), kind: keyGamepadLeftStick, name: "gamepad_lstick_down"}
	KeyGamepadLStickLeft  = Key{code: int(stickLeft), kind: keyGamepadLeftStick, name: "gamepad_lstick_left"}
	KeyGamepadRStickUp    = Key{code: int(stickUp), kind: keyGamepadRightStick, name: "gamepad_rstick_up"}
	KeyGamepadRStickRight = Key{code: int(stickRight), kind: keyGamepadRightStick, name: "gamepad_rstick_right"}
	KeyGamepadRStickDown  = Key{code: int(stickDown), kind: keyGamepadRightStick, name: "gamepad_rstick_down"}
	KeyGamepadRStickLeft  = Key{code: int(stickLeft), kind: keyGamepadRightStick, name: "gamepad_rstick_left"}

	KeyGamepadA Key = Key{code: int(ebiten.StandardGamepadButtonRightBottom), kind: keyGamepad, name: "gamepad_a"}
	KeyGamepadB Key = Key{code: int(ebiten.StandardGamepadButtonRightRight), kind: keyGamepad, name: "gamepad_b"}
	KeyGamepadX Key = Key{code: int(ebiten.StandardGamepadButtonRightLeft), kind: keyGamepad, name: "gamepad_x"}
	KeyGamepadY Key = Key{code: int(ebiten.StandardGamepadButtonRightTop), kind: keyGamepad, name: "gamepad_y"}

	KeyGamepadL1 Key = Key{code: int(ebiten.StandardGamepadButtonFrontTopLeft), kind: keyGamepad, name: "gamepad_l1"}
	KeyGamepadL2 Key = Key{code: int(ebiten.StandardGamepadButtonFrontBottomLeft), kind: keyGamepad, name: "gamepad_l2"}
	KeyGamepadR1 Key = Key{code: int(ebiten.StandardGamepadButtonFrontTopRight), kind: keyGamepad, name: "gamepad_r1"}
	KeyGamepadR2 Key = Key{code: int(ebiten.StandardGamepadButtonFrontBottomRight), kind: keyGamepad, name: "gamepad_r2"}
)
