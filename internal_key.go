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
	keyMouse
	keyMouseWithCtrl
	keyMouseWithShift
	keyMouseWithCtrlShift
	keyTouch
	keyWheel
)

type touchCode int

const (
	touchUnknown touchCode = iota
	touchTap
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

const (
	// Make "unused" linter happy.
	// We have these unknown values as a programming habbit of having
	// an invalid (undefined) state of the enum value.
	_ = touchUnknown
	_ = wheelUnknown
	_ = stickUnknown
)
