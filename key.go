package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ebitengine-input/keyname"
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
	KeyMouseLeft   = Key{code: int(ebiten.MouseButtonLeft), kind: keyMouse, name: keyname.MouseLeft}
	KeyMouseRight  = Key{code: int(ebiten.MouseButtonRight), kind: keyMouse, name: keyname.MouseRight}
	KeyMouseMiddle = Key{code: int(ebiten.MouseButtonMiddle), kind: keyMouse, name: keyname.MouseMiddle}
)

// Touch keys.
var (
	KeyTouchTap = Key{code: int(touchTap), kind: keyTouch, name: keyname.ScreenTap}
)

// Keyboard keys.
var (
	KeyLeft  Key = Key{code: int(ebiten.KeyLeft), name: keyname.KeyboardLeft}
	KeyRight Key = Key{code: int(ebiten.KeyRight), name: keyname.KeyboardRight}
	KeyUp    Key = Key{code: int(ebiten.KeyUp), name: keyname.KeyboardUp}
	KeyDown  Key = Key{code: int(ebiten.KeyDown), name: keyname.KeyboardDown}

	KeyTab Key = Key{code: int(ebiten.KeyTab), name: keyname.KeyboardTab}

	Key0 Key = Key{code: int(ebiten.Key0), name: keyname.Keyboard0}
	Key1 Key = Key{code: int(ebiten.Key1), name: keyname.Keyboard1}
	Key2 Key = Key{code: int(ebiten.Key2), name: keyname.Keyboard2}
	Key3 Key = Key{code: int(ebiten.Key3), name: keyname.Keyboard3}
	Key4 Key = Key{code: int(ebiten.Key4), name: keyname.Keyboard4}
	Key5 Key = Key{code: int(ebiten.Key5), name: keyname.Keyboard5}
	Key6 Key = Key{code: int(ebiten.Key6), name: keyname.Keyboard6}
	Key7 Key = Key{code: int(ebiten.Key7), name: keyname.Keyboard7}
	Key8 Key = Key{code: int(ebiten.Key8), name: keyname.Keyboard8}
	Key9 Key = Key{code: int(ebiten.Key9), name: keyname.Keyboard9}

	KeyQ Key = Key{code: int(ebiten.KeyQ), name: keyname.KeyboardQ}
	KeyW Key = Key{code: int(ebiten.KeyW), name: keyname.KeyboardW}
	KeyE Key = Key{code: int(ebiten.KeyE), name: keyname.KeyboardE}
	KeyR Key = Key{code: int(ebiten.KeyR), name: keyname.KeyboardR}
	KeyT Key = Key{code: int(ebiten.KeyT), name: keyname.KeyboardT}
	KeyY Key = Key{code: int(ebiten.KeyY), name: keyname.KeyboardY}
	KeyU Key = Key{code: int(ebiten.KeyY), name: keyname.KeyboardU}
	KeyI Key = Key{code: int(ebiten.KeyY), name: keyname.KeyboardI}
	KeyO Key = Key{code: int(ebiten.KeyY), name: keyname.KeyboardO}
	KeyP Key = Key{code: int(ebiten.KeyY), name: keyname.KeyboardP}
	KeyA Key = Key{code: int(ebiten.KeyA), name: keyname.KeyboardA}
	KeyS Key = Key{code: int(ebiten.KeyS), name: keyname.KeyboardS}
	KeyD Key = Key{code: int(ebiten.KeyD), name: keyname.KeyboardD}
	KeyF Key = Key{code: int(ebiten.KeyD), name: keyname.KeyboardF}
	KeyG Key = Key{code: int(ebiten.KeyD), name: keyname.KeyboardG}
	KeyH Key = Key{code: int(ebiten.KeyD), name: keyname.KeyboardH}
	KeyJ Key = Key{code: int(ebiten.KeyD), name: keyname.KeyboardJ}
	KeyK Key = Key{code: int(ebiten.KeyD), name: keyname.KeyboardK}
	KeyL Key = Key{code: int(ebiten.KeyD), name: keyname.KeyboardL}
	KeyZ Key = Key{code: int(ebiten.KeyD), name: keyname.KeyboardZ}
	KeyX Key = Key{code: int(ebiten.KeyD), name: keyname.KeyboardX}
	KeyC Key = Key{code: int(ebiten.KeyD), name: keyname.KeyboardC}
	KeyV Key = Key{code: int(ebiten.KeyD), name: keyname.KeyboardV}
	KeyB Key = Key{code: int(ebiten.KeyD), name: keyname.KeyboardB}
	KeyN Key = Key{code: int(ebiten.KeyD), name: keyname.KeyboardN}
	KeyM Key = Key{code: int(ebiten.KeyD), name: keyname.KeyboardM}

	KeyEscape Key = Key{code: int(ebiten.KeyEscape), name: keyname.KeyboardEscape}
	KeyEnter  Key = Key{code: int(ebiten.KeyEnter), name: keyname.KeyboardEnter}

	KeySpace Key = Key{code: int(ebiten.KeySpace), name: keyname.KeyboardSpace}
)

// Gamepad keys.
var (
	KeyGamepadStart  Key = Key{code: int(ebiten.StandardGamepadButtonCenterRight), kind: keyGamepad, name: keyname.GamepadStart}
	KeyGamepadSelect Key = Key{code: int(ebiten.StandardGamepadButtonCenterLeft), kind: keyGamepad, name: keyname.GamepadSelect}
	KeyGamepadMiddle Key = Key{code: int(ebiten.StandardGamepadButtonCenterCenter), kind: keyGamepad, name: keyname.GamepadMiddle}

	KeyGamepadUp    Key = Key{code: int(ebiten.StandardGamepadButtonLeftTop), kind: keyGamepad, name: keyname.GamepadUp}
	KeyGamepadRight Key = Key{code: int(ebiten.StandardGamepadButtonLeftRight), kind: keyGamepad, name: keyname.GamepadRight}
	KeyGamepadDown  Key = Key{code: int(ebiten.StandardGamepadButtonLeftBottom), kind: keyGamepad, name: keyname.GamepadDown}
	KeyGamepadLeft  Key = Key{code: int(ebiten.StandardGamepadButtonLeftLeft), kind: keyGamepad, name: keyname.GamepadLeft}

	KeyGamepadLStickUp    = Key{code: int(stickUp), kind: keyGamepadLeftStick, name: keyname.GamepadLStickUp}
	KeyGamepadLStickRight = Key{code: int(stickRight), kind: keyGamepadLeftStick, name: keyname.GamepadLStickRight}
	KeyGamepadLStickDown  = Key{code: int(stickDown), kind: keyGamepadLeftStick, name: keyname.GamepadLStickDown}
	KeyGamepadLStickLeft  = Key{code: int(stickLeft), kind: keyGamepadLeftStick, name: keyname.GamepadLStickLeft}
	KeyGamepadRStickUp    = Key{code: int(stickUp), kind: keyGamepadRightStick, name: keyname.GamepadRStickUp}
	KeyGamepadRStickRight = Key{code: int(stickRight), kind: keyGamepadRightStick, name: keyname.GamepadRStickRight}
	KeyGamepadRStickDown  = Key{code: int(stickDown), kind: keyGamepadRightStick, name: keyname.GamepadRStickDown}
	KeyGamepadRStickLeft  = Key{code: int(stickLeft), kind: keyGamepadRightStick, name: keyname.GamepadRStickLeft}

	KeyGamepadA Key = Key{code: int(ebiten.StandardGamepadButtonRightBottom), kind: keyGamepad, name: keyname.GamepadA}
	KeyGamepadB Key = Key{code: int(ebiten.StandardGamepadButtonRightRight), kind: keyGamepad, name: keyname.GamepadB}
	KeyGamepadX Key = Key{code: int(ebiten.StandardGamepadButtonRightLeft), kind: keyGamepad, name: keyname.GamepadX}
	KeyGamepadY Key = Key{code: int(ebiten.StandardGamepadButtonRightTop), kind: keyGamepad, name: keyname.GamepadY}

	KeyGamepadL1 Key = Key{code: int(ebiten.StandardGamepadButtonFrontTopLeft), kind: keyGamepad, name: keyname.GamepadL1}
	KeyGamepadL2 Key = Key{code: int(ebiten.StandardGamepadButtonFrontBottomLeft), kind: keyGamepad, name: keyname.GamepadL2}
	KeyGamepadR1 Key = Key{code: int(ebiten.StandardGamepadButtonFrontTopRight), kind: keyGamepad, name: keyname.GamepadR1}
	KeyGamepadR2 Key = Key{code: int(ebiten.StandardGamepadButtonFrontBottomRight), kind: keyGamepad, name: keyname.GamepadR2}
)
