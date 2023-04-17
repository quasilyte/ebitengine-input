package input

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

type gamepadModel int

const (
	gamepadUnknown gamepadModel = iota
	gamepadStandard
	gamepadFirefoxXinput
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

	axisCount      int
	axisValues     [8]float64
	prevAxisValues [8]float64
}

func isDPadButton(code int) bool {
	switch ebiten.StandardGamepadButton(code) {
	case ebiten.StandardGamepadButtonLeftTop:
		return true
	case ebiten.StandardGamepadButtonLeftRight:
		return true
	case ebiten.StandardGamepadButtonLeftBottom:
		return true
	case ebiten.StandardGamepadButtonLeftLeft:
		return true
	default:
		return false
	}
}

func firefoxXinputToXbox(b ebiten.StandardGamepadButton) ebiten.GamepadButton {
	switch b {
	case ebiten.StandardGamepadButtonCenterLeft:
		return 6
	case ebiten.StandardGamepadButtonCenterRight:
		return 7
	case ebiten.StandardGamepadButtonRightStick:
		return 10
	case ebiten.StandardGamepadButtonLeftStick:
		return 9
	default:
		return ebiten.GamepadButton(b)
	}
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
