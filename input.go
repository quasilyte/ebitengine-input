package input

import (
	"strings"
)

// Action is an ID that represents an abstract action that can be activeted by the input.
type Action uint32

// Keymap associates a list of keys with an action.
// Any of the keys from the list can activate the action.
type Keymap map[Action][]Key

// Clone creates a deep copy of a keymap.
// The returned keymap can be modified without changing the original keymap.
func (m Keymap) Clone() Keymap {
	cloned := make(Keymap, len(m))
	for k, list := range m {
		clonedList := make([]Key, len(list))
		copy(clonedList, list)
		cloned[k] = clonedList
	}
	return cloned
}

// InputDeviceKind is used as a bit mask to select the enabled input devices.
// See constants like KeyboardInput and GamepadInput.
// Combine them like KeyboardInput|GamepadInput to get a bit mask that includes multiple entries.
// Use AnyInput if you want to have a mask covering all devices.
type InputDeviceKind uint8

const (
	KeyboardInput InputDeviceKind = 1 << iota
	GamepadInput
	MouseInput
	TouchInput
)

// String returns a pretty-printed representation of the input device mask.
func (d InputDeviceKind) String() string {
	if d == 0 {
		return "<empty>"
	}
	parts := make([]string, 0, 4)
	if d&KeyboardInput != 0 {
		parts = append(parts, "keyboard")
	}
	if d&GamepadInput != 0 {
		parts = append(parts, "gamepad")
	}
	if d&MouseInput != 0 {
		parts = append(parts, "mouse")
	}
	if d&TouchInput != 0 {
		parts = append(parts, "touch")
	}
	if len(parts) == 0 {
		return "<invalid>"
	}
	return strings.Join(parts, "|")
}

// AnyInput includes all input devices.
const AnyInput InputDeviceKind = KeyboardInput | GamepadInput | MouseInput | TouchInput
