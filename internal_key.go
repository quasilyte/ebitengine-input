package input

type keyKind uint8

const (
	keyKeyboard keyKind = iota
	keyKeyboardWithCtrl
	keyKeyboardWithShift
	keyGamepad
	keyGamepadLeftStick
	keyGamepadRightStick
	keyMouse
	keyMouseWithCtrl
	keyMouseWithShift
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
