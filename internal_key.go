package input

type keyKind uint8

const (
	keyKeyboard keyKind = iota
	keyKeyboardWithCtrl
	keyKeyboardWithShift
	keyKeyboardWithCtrlShift
	keyGamepad
	keyGamepadLeftStick
	keyGamepadRightStick
	keyGamepadStickMotion
	keyMouse
	keyMouseWithCtrl
	keyMouseWithShift
	keyMouseWithCtrlShift
	keyTouch
	keyTouchDrag
	keyWheel
	keySimulated
)

type touchCode int

const (
	touchUnknown touchCode = iota
	touchTap
	touchLongTap
	touchDrag
)

type wheelCode int

const (
	wheelUnknown wheelCode = iota
	wheelUp
	wheelDown
	wheelVertical
)

type stickCode int

const (
	stickUnknown stickCode = iota
	stickUp
	stickRight
	stickDown
	stickLeft
)

type keyKindFlag uint8

const (
	keyFlagHasPos keyKindFlag = 1 << iota
	keyFlagNeedID
)

func keyHasPos(k keyKind) bool { return keyKindFlagTable[k]&keyFlagHasPos != 0 }
func keyNeedID(k keyKind) bool { return keyKindFlagTable[k]&keyFlagNeedID != 0 }

// Using a 256-byte LUT to get a fast map-like lookup without a bound check.
var keyKindFlagTable = [256]keyKindFlag{
	keySimulated: keyFlagHasPos | keyFlagNeedID,

	keyKeyboard:              0,
	keyKeyboardWithCtrl:      0,
	keyKeyboardWithShift:     0,
	keyKeyboardWithCtrlShift: 0,

	keyGamepad:           keyFlagNeedID,
	keyGamepadLeftStick:  keyFlagNeedID,
	keyGamepadRightStick: keyFlagNeedID,

	keyGamepadStickMotion: keyFlagHasPos | keyFlagNeedID,

	keyMouse:              keyFlagHasPos,
	keyMouseWithCtrl:      keyFlagHasPos,
	keyMouseWithShift:     keyFlagHasPos,
	keyMouseWithCtrlShift: keyFlagHasPos,
	keyTouch:              keyFlagHasPos,
	keyWheel:              keyFlagHasPos,
}
