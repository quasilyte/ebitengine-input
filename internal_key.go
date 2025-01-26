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
	keyMouseDrag
	keyTouch
	keyTouchDrag
	keyWheel
	keyWheelWithCtrl
	keyWheelWithShift
	keyWheelWithCtrlShift
	keySimulated
)

func (k keyKind) device() DeviceKind {
	switch k {
	case keyKeyboard, keyKeyboardWithCtrl, keyKeyboardWithShift, keyKeyboardWithCtrlShift:
		return KeyboardDevice
	case keyGamepad, keyGamepadLeftStick, keyGamepadRightStick, keyGamepadStickMotion:
		return GamepadDevice
	case keyMouse, keyMouseDrag:
		return MouseDevice
	case keyWheel, keyWheelWithCtrl, keyWheelWithShift, keyWheelWithCtrlShift:
		return MouseDevice
	case keyMouseWithCtrl, keyMouseWithShift, keyMouseWithCtrlShift:
		return MouseDevice | KeyboardDevice
	case keyTouch, keyTouchDrag:
		return TouchDevice
	default:
		return KeyboardDevice
	}
}

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
	keyFlagHasDuration
)

func keyHasPos(k keyKind) bool      { return keyKindFlagTable[k]&keyFlagHasPos != 0 }
func keyNeedID(k keyKind) bool      { return keyKindFlagTable[k]&keyFlagNeedID != 0 }
func keyHasDuration(k keyKind) bool { return keyKindFlagTable[k]&keyFlagHasDuration != 0 }

// Using a 256-byte LUT to get a fast map-like lookup without a bound check.
var keyKindFlagTable = [256]keyKindFlag{
	keySimulated: keyFlagHasPos | keyFlagNeedID,

	keyKeyboard:              keyFlagHasDuration,
	keyKeyboardWithCtrl:      keyFlagHasDuration,
	keyKeyboardWithShift:     keyFlagHasDuration,
	keyKeyboardWithCtrlShift: keyFlagHasDuration,

	keyGamepad:           keyFlagNeedID,
	keyGamepadLeftStick:  keyFlagNeedID,
	keyGamepadRightStick: keyFlagNeedID,

	keyGamepadStickMotion: keyFlagHasPos | keyFlagNeedID,

	keyMouse:              keyFlagHasPos,
	keyMouseWithCtrl:      keyFlagHasPos,
	keyMouseWithShift:     keyFlagHasPos,
	keyMouseWithCtrlShift: keyFlagHasPos,
	keyMouseDrag:          keyFlagHasPos,
	keyTouch:              keyFlagHasPos,
	keyWheel:              keyFlagHasPos,
	keyWheelWithCtrl:      keyFlagHasPos,
	keyWheelWithShift:     keyFlagHasPos,
	keyWheelWithCtrlShift: keyFlagHasPos,
}
