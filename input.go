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

// uniqueKeys returns a list of [Key]s that do not repeat.
func uniqueKeys(keys []Key) (unique []Key) {
	known := make(map[Key]struct{})
	var ok bool
	for _, key := range keys {
		if _, ok = known[key]; ok {
			continue // skip
		}
		known[key] = struct{}{}
		unique = append(unique, key)
	}
	return unique
}

// MergeKeymaps merges a list of [Keymap]s into one.
// Given maps are not modified.
// Resulting map contains no references to given maps.
func MergeKeymaps(maps ...Keymap) Keymap {
	merged := make(Keymap)
	for _, m := range maps {
		for k, list := range m {
			merged[k] = uniqueKeys(append(merged[k], list...))
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
