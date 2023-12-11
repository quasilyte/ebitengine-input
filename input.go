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

// MergeKeymaps merges a list of keymaps into one.
//
// Given maps are not modified.
// The resulting keymap contains no references to input keymaps.
//
// Any key duplicates for a single action are ignored, therefore:
//
//	{Action1: [KeyA, KeyB]} + {Action1: [KeyC, KeyB]} = {Action1: [KeyA, KeyB, KeyC]}
//
// The keys order depends on the input keymaps arguments order.
// Keymaps with keys of higher priorities should go first.
// So, if you have a keyboardKeymap and gamepadKeymap keymaps,
// passing gamepadKeymap before keyboardKeymap would make the gamepad key
// binds take priority: MergeKeymaps(gamepadKeymap, keyboardKeymap).
func MergeKeymaps(maps ...Keymap) Keymap {
	type keyCombination struct {
		a Action
		k Key
	}
	keyCombinations := make(map[keyCombination]struct{})
	merged := make(Keymap)
	for _, m := range maps {
		for a, keys := range m {
			for _, k := range keys {
				ck := keyCombination{a, k}
				if _, ok := keyCombinations[ck]; ok {
					continue
				}
				keyCombinations[ck] = struct{}{}
				merged[a] = append(merged[a], k)
			}
		}
	}
	return merged
}

// DeviceKind is used as a bit mask to select the enabled input devices.
// See constants like KeyboardInput and GamepadInput.
// Combine them like KeyboardInput|GamepadInput to get a bit mask that includes multiple entries.
// Use AnyDevice if you want to have a mask covering all devices.
type DeviceKind uint8

const (
	KeyboardDevice DeviceKind = 1 << iota
	GamepadDevice
	MouseDevice
	TouchDevice
)

// String returns a pretty-printed representation of the input device mask.
func (d DeviceKind) String() string {
	if d == 0 {
		return "<empty>"
	}
	parts := make([]string, 0, 4)
	if d&KeyboardDevice != 0 {
		parts = append(parts, "keyboard")
	}
	if d&GamepadDevice != 0 {
		parts = append(parts, "gamepad")
	}
	if d&MouseDevice != 0 {
		parts = append(parts, "mouse")
	}
	if d&TouchDevice != 0 {
		parts = append(parts, "touch")
	}
	if len(parts) == 0 {
		return "<invalid>"
	}
	return strings.Join(parts, "|")
}

// AnyDevice includes all input devices.
const AnyDevice DeviceKind = KeyboardDevice | GamepadDevice | MouseDevice | TouchDevice
