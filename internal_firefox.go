//go:build !js

package input

func isFirefox() bool {
	return false
}

func guessFirefoxGamepadModel(id int) gamepadModel {
	panic("should not be called")
}
