package input

import "testing"

func TestUniqueKeys(t *testing.T) {
	unique := uniqueKeys([]Key{KeyUp, KeyUp, KeyUp})
	if len(unique) != 1 {
		t.Fatal("duplicates were not removed")
	}
	if unique[0] != KeyUp {
		t.Fatal("unexpected key")
	}
}

func TestKeymapMerge(t *testing.T) {
	keyboardKeyMap := Keymap{
		3: {KeyDown, KeyS},
		4: {KeySpace},
		5: {KeyShift},
	}

	mouseKeymap := Keymap{
		4: {KeyMouseLeft, KeySpace}, // extra duplicate
		5: {KeyMouseRight},
	}

	merged := MergeKeymaps(keyboardKeyMap, mouseKeymap)

	if l := len(merged); l != 3 {
		t.Fatalf("key map contains%d elements instead of expected 3", l)
	}

	if l := len(merged[4]); l != 2 {
		t.Fatalf("duplicate was not removed: %+v", merged[4])
	}
}
