package input

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestScanKey(t *testing.T) {
	tests := []struct {
		keys []ebiten.Key
		want Key
	}{
		// Sanity tests.
		{[]ebiten.Key{}, Key{}},

		// The simple cases with a single key.
		{[]ebiten.Key{ebiten.KeyB}, KeyB},
		{[]ebiten.Key{ebiten.KeyEnter}, KeyEnter},
		{[]ebiten.Key{ebiten.KeyControlLeft, ebiten.KeyControl}, KeyControlLeft},
		{[]ebiten.Key{ebiten.KeyControlRight, ebiten.KeyControl}, KeyControlRight},
		{[]ebiten.Key{ebiten.KeyControl, ebiten.KeyControlRight}, KeyControlRight},

		// Multiple key candidates without a way to merge them into a single Key.
		{[]ebiten.Key{ebiten.KeyB, ebiten.KeyA}, KeyA},
		{[]ebiten.Key{ebiten.KeyA, ebiten.KeyB}, KeyA},

		// Control modifiers.
		{[]ebiten.Key{ebiten.KeyC, ebiten.KeyControlLeft, ebiten.KeyControl}, KeyWithModifier(KeyC, ModControl)},
		{[]ebiten.Key{ebiten.KeyC, ebiten.KeyControl, ebiten.KeyControlLeft}, KeyWithModifier(KeyC, ModControl)},
		{[]ebiten.Key{ebiten.KeyControl, ebiten.KeyControlLeft, ebiten.KeyC}, KeyWithModifier(KeyC, ModControl)},
		{[]ebiten.Key{ebiten.KeyE, ebiten.KeyControlLeft, ebiten.KeyControl}, KeyWithModifier(KeyE, ModControl)},
		{[]ebiten.Key{ebiten.KeyControlLeft, ebiten.KeyControl, ebiten.KeyE}, KeyWithModifier(KeyE, ModControl)},

		// Shift modifiers.
		{[]ebiten.Key{ebiten.KeyF, ebiten.KeyShiftLeft, ebiten.KeyShift}, KeyWithModifier(KeyF, ModShift)},

		// Control+Shift modifiers.
		{[]ebiten.Key{ebiten.KeyC, ebiten.KeyControlLeft, ebiten.KeyControl, ebiten.KeyShiftLeft, ebiten.KeyShift}, KeyWithModifier(KeyC, ModControlShift)},
		{[]ebiten.Key{ebiten.KeyA, ebiten.KeyControlLeft, ebiten.KeyControl, ebiten.KeyShiftLeft, ebiten.KeyShift}, KeyWithModifier(KeyA, ModControlShift)},
		{[]ebiten.Key{ebiten.KeyA, ebiten.KeyControlLeft, ebiten.KeyShiftLeft}, KeyWithModifier(KeyA, ModControlShift)},
		{[]ebiten.Key{ebiten.KeyA, ebiten.KeyControlRight, ebiten.KeyShiftRight}, KeyWithModifier(KeyA, ModControlShift)},
		{[]ebiten.Key{ebiten.KeyControlLeft, ebiten.KeyA, ebiten.KeyShiftRight}, KeyWithModifier(KeyA, ModControlShift)},
	}

	for i, test := range tests {
		have, _ := scanKey(test.keys)
		if have != test.want {
			t.Fatalf("test[%d] failed:\nhave: %s (%#v)\nwant: %s (%#v)",
				i, have, have, test.want, test.want)
		}
	}
}
