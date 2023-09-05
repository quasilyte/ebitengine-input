package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// KeyScanStatus represents the KeyScanner.Scan operation result.
type KeyScanStatus int

const (
	KeyScanUnchanged KeyScanStatus = iota
	KeyScanChanged
	KeyScanCompleted
)

// KeyScanner checks the currently pressed keys and buttons and tries to map them
// to a local Key type that can be used in a Keymap.
//
// Use NewKeyScanner to create a usable object of this type.
//
// Experimental: this is a part of a key remapping API, which is not stable yet.
type KeyScanner struct {
	lastNumKeys int
	canScan     bool
	key         Key
	h           *Handler
}

// NewKeyScanner creates a key scanner for the specifier input Handler.
//
// You don't have to create a new scanner for every remap; they can be reused.
//
// It's important to have the correct Handler though: their ID is used to
// check the appropriate device keys.
//
// Experimental: this is a part of a key remapping API, which is not stable yet.
func NewKeyScanner(h *Handler) *KeyScanner {
	return &KeyScanner{h: h}
}

// Scan reads the buttons state and tries to map them to a Key.
//
// It's intended to work with keyboard keys as well as gamepad buttons,
// but right now it only works for the keyboard.
//
// This function should be called on every frame where you're reading
// the new keybind combination.
// See the remap example for more info.
//
// The function can return these result statuses:
// * Unchanged - nothing updated since the last Scan() operation
// * Changed - some keys changed, you may want to update the prompt to the user
// * Completed - the user finished specifying the keys combination, you can use the Key as a new binding
func (s *KeyScanner) Scan() (Key, KeyScanStatus) {
	// TODO: respect the enabled input devices.
	// TODO: scan the gamepad buttons as well.

	// Note that this function may not be needed by some users,
	// so we're better of making it as independent as possible, so it
	// doesn't make the package more expensive if you don't use it.
	//
	// This function doesn't have to be very fast, but it should be relatively
	// inexpensive for the "no keys were pressed" case.
	// When some keys combo is being pressed, it's OK to spend some resources.

	// This slice is stack-allocated; for the most cases, 4 keys are enough.
	keys := make([]ebiten.Key, 0, 4)
	keys = inpututil.AppendPressedKeys(keys)

	if !s.canScan {
		if len(keys) != 0 {
			return Key{}, KeyScanUnchanged
		}
		s.canScan = true
	}

	if len(keys) == s.lastNumKeys {
		// It's either empty or we're still collecting the keys.
		return Key{}, KeyScanUnchanged
	}

	if len(keys) < s.lastNumKeys {
		// One or more keys are released.
		// Consider it to be a confirmation event.
		result := s.key
		s.lastNumKeys = 0
		s.key = Key{}
		s.canScan = false
		return result, KeyScanCompleted
	}

	s.lastNumKeys = len(keys)

	k, ok := scanKey(keys)
	status := KeyScanUnchanged
	if ok {
		s.key = k
		status = KeyScanChanged
	}
	return k, status
}

func scanKey(keys []ebiten.Key) (Key, bool) {
	if len(keys) == 0 {
		return Key{}, false
	}

	containsKeyCode := func(keys []ebiten.Key, code int) bool {
		for _, k := range keys {
			if int(k) == code {
				return true
			}
		}
		return false
	}

	// Parse the keys combination into something that this library can handle.

	// Round 1: walk the actual keys that are being pressed and collect the modifiers.
	// Remove the modifiers from the slice (inplace).
	var ctrlKey Key
	var shiftKey Key
	keysWithoutMods := keys[:0]
	for _, k := range keys {
		switch k {
		case ebiten.KeyControl, ebiten.KeyShift:
			// Just omit them from the slice.
		case ebiten.KeyControlLeft:
			ctrlKey = KeyControlLeft
		case ebiten.KeyControlRight:
			ctrlKey = KeyControlRight
		case ebiten.KeyShiftLeft:
			shiftKey = KeyShiftLeft
		case ebiten.KeyShiftRight:
			shiftKey = KeyShiftRight
		default:
			keysWithoutMods = append(keysWithoutMods, k)
		}
	}
	hasCtrl := ctrlKey.name != ""
	hasShift := shiftKey.name != ""

	var mappedKey Key

	// Round 2: map the Ebitengine keys to the local types.
	// In theory, we could generate a big LUT to make this mapping very fast.
	// But this would mean more data reserved for this package.
	// Since this part of the code is not that performance-sensitive,
	// we'll handle it in a less efficient, but less memory-hungry way.
Loop:
	for _, lk := range allKeys {
		switch lk.kind {
		case keyKeyboard:
			if containsKeyCode(keysWithoutMods, lk.code) {
				mappedKey = lk
				break Loop
			}
		}
	}

	if mappedKey.name == "" {
		switch {
		case hasCtrl:
			return ctrlKey, true
		case hasShift:
			return shiftKey, true
		}
	}

	var keymod KeyModifier
	switch {
	case hasCtrl && hasShift:
		keymod = ModControlShift
	case hasCtrl:
		keymod = ModControl
	case hasShift:
		keymod = ModShift
	}
	if keymod != ModUnknown {
		switch mappedKey.kind {
		case keyKeyboard, keyMouse:
			mappedKey = KeyWithModifier(mappedKey, keymod)
		}
	}

	return mappedKey, mappedKey.name != ""
}
