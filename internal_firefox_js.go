package input

import (
	"strings"
	"syscall/js"
)

func isFirefox() bool {
	ua := js.Global().Get("navigator").Get("userAgent").String()
	return strings.Contains(strings.ToLower(ua), "firefox")
}

func guessFirefoxGamepadModel(id int) gamepadModel {
	gamepads := js.Global().Get("navigator").Call("getGamepads")
	if gamepads.IsNull() || gamepads.Type() != js.TypeObject {
		return gamepadUnknown
	}
	g := gamepads.Index(id)
	if g.IsNull() {
		return gamepadUnknown
	}
	gamepadID := strings.ToLower(g.Get("id").String())
	for _, pattern := range firefoxKnownXinput {
		if strings.Contains(gamepadID, pattern) {
			return gamepadFirefoxXinput
		}
	}
	return gamepadUnknown
}

var firefoxKnownXinput = []string{
	// Generic keys.
	"xinput",
	"x-input",
	"x_input",
	"xbox",
	"x-box",
	"x_box",

	// Specific models that do not contain any xinput keys in their name.
	"logitech gamepad f310",
}
