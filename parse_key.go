package input

import (
	"errors"
	"sort"
	"strings"
)

// ParseKeys tries to construct an appropriate Key object given its name.
//
// It can also be used as a string->key constructor:
//
//	ParseKey("left")         // returns KeyLeft
//	ParseKey("gamepad_left") // returns KeyGamepadLeft
//
// The format is one of the following:
//
//   - keyname
//   - mod+keyname
//   - mod+mod+keyname
//
// Some valid input examples:
//
//   - "gamepad_left"
//   - "left"
//   - "ctrl+left"
//   - "ctrl+shift+left"
//   - "shift+ctrl+left"
//
// See Handler.ActionKeyNames() for more information about the key names.
func ParseKey(s string) (Key, error) {
	plusPos := strings.LastIndex(s, "+")
	if plusPos == -1 {
		k := keyByName(s)
		if (k == Key{}) {
			return k, errors.New("unknown key: " + s)
		}
		return k, nil
	}
	modName := s[:plusPos]
	keyName := s[plusPos+1:]
	mod := keyModifierByName(modName)
	if mod == ModUnknown {
		return Key{}, errors.New("unknown key modifier: " + modName)
	}
	k := keyByName(keyName)
	if (k == Key{}) {
		return k, errors.New("unknown key: " + keyName)
	}
	return KeyWithModifier(k, mod), nil
}

func keyModifierByName(name string) KeyModifier {
	switch name {
	case "ctrl":
		return ModControl
	case "shift":
		return ModShift
	case "ctrl+shift", "shift+ctrl":
		return ModControlShift
	default:
		return ModUnknown
	}
}

func keyByName(name string) Key {
	// Keys are sorted by a name, so we can use a binary search here.
	i := sort.Search(len(allKeys), func(i int) bool {
		return allKeys[i].name >= name
	})
	if i < len(allKeys) && allKeys[i].name == name {
		return allKeys[i]
	}
	return Key{}
}
