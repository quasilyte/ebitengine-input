package input

import (
	"fmt"
	"reflect"
	"testing"
)

func TestKeymapMerge(t *testing.T) {
	tests := []struct {
		keymaps []Keymap
		want    Keymap
	}{
		// A simple case with 4 keymaps with no duplicates.
		{
			[]Keymap{
				{
					3: {KeyDown, KeyS},
					4: {KeySpace},
					5: {KeyShift},
				},
				{
					4: {KeyMouseLeft},
					5: {KeyMouseRight},
				},
				{
					6: {KeyMouseRight},
				},
				{
					7: {KeyGamepadA, KeyGamepadB},
					3: {KeyGamepadA, KeyGamepadB},
				},
				{
					3: {KeyGamepadL1},
					6: {KeyGamepadL2},
				},
			},
			Keymap{
				3: {KeyDown, KeyS, KeyGamepadA, KeyGamepadB, KeyGamepadL1},
				4: {KeySpace, KeyMouseLeft},
				5: {KeyShift, KeyMouseRight},
				6: {KeyMouseRight, KeyGamepadL2},
				7: {KeyGamepadA, KeyGamepadB},
			},
		},

		{
			[]Keymap{
				{
					3: {KeyDown, KeyS},
					4: {KeySpace},
					5: {KeyShift},
				},
				{
					4: {KeyMouseLeft, KeySpace}, // extra duplicate
					5: {KeyMouseRight},
				},
			},
			Keymap{
				3: {KeyDown, KeyS},
				4: {KeySpace, KeyMouseLeft},
				5: {KeyShift, KeyMouseRight},
			},
		},

		// Merging with 3 keymaps, checking that the priority is preserved.
		{
			[]Keymap{
				{
					3: {KeyGamepadA},
					4: {KeySpace, KeyMouseLeft},
				},
				{
					3: {KeyDown, KeyS},
					4: {KeySpace},
					5: {KeyShift},
				},
				{
					4: {KeyMouseLeft, KeySpace},
					5: {KeyMouseRight},
				},
			},
			Keymap{
				3: {KeyGamepadA, KeyDown, KeyS},
				4: {KeySpace, KeyMouseLeft},
				5: {KeyShift, KeyMouseRight},
			},
		},

		// Merging with an empty keymap.
		{
			[]Keymap{
				{
					3: {KeyDown, KeyS},
					4: {KeySpace},
					5: {KeyShift},
				},
				{},
			},
			Keymap{
				3: {KeyDown, KeyS},
				4: {KeySpace},
				5: {KeyShift},
			},
		},

		// Merging identical keymaps.
		{
			[]Keymap{
				{
					3: {KeyDown, KeyS},
					4: {KeySpace},
					5: {KeyShift},
				},
				{
					3: {KeyDown, KeyS},
					4: {KeySpace},
					5: {KeyShift},
				},
				{
					3: {KeyDown, KeyS},
					4: {KeySpace},
					5: {KeyShift},
				},
			},
			Keymap{
				3: {KeyDown, KeyS},
				4: {KeySpace},
				5: {KeyShift},
			},
		},

		// Merging a single map results in the same keymap.
		{
			[]Keymap{
				{
					3: {KeyDown, KeyS},
					4: {KeySpace},
					5: {KeyShift},
				},
			},
			Keymap{
				3: {KeyDown, KeyS},
				4: {KeySpace},
				5: {KeyShift},
			},
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			have := MergeKeymaps(test.keymaps...)
			if !reflect.DeepEqual(have, test.want) {
				t.Fatalf("invalid merge results:\nhave: %#v\nwant: %#v\ninputs: #%v", have, test.want, test.keymaps)
			}
		})
	}
}
