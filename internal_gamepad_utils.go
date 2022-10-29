package input

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

type gamepadModel int

const (
	gamepadUnknown gamepadModel = iota
	gamepadStandard
	gamepadMicront
)

func guessGamepadModel(s string) gamepadModel {
	s = strings.ToLower(s)
	if s == "micront" {
		return gamepadMicront
	}
	return gamepadUnknown
}

type gamepadInfo struct {
	model     gamepadModel
	modelName string
}

func microntToXbox(b ebiten.StandardGamepadButton) ebiten.GamepadButton {
	switch b {
	case ebiten.StandardGamepadButtonLeftTop:
		return ebiten.GamepadButton12
	case ebiten.StandardGamepadButtonLeftRight:
		return ebiten.GamepadButton13
	case ebiten.StandardGamepadButtonLeftBottom:
		return ebiten.GamepadButton14
	case ebiten.StandardGamepadButtonLeftLeft:
		return ebiten.GamepadButton15

	case ebiten.StandardGamepadButtonRightTop:
		return ebiten.GamepadButton0
	case ebiten.StandardGamepadButtonRightRight:
		return ebiten.GamepadButton1
	case ebiten.StandardGamepadButtonRightBottom:
		return ebiten.GamepadButton2
	case ebiten.StandardGamepadButtonRightLeft:
		return ebiten.GamepadButton3

	default:
		return ebiten.GamepadButton(b)
	}
}
