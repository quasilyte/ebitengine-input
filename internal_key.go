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
	keyTouch
)

type touchCode int

const (
	touchUnknown touchCode = iota
	touchTap
)

type stickCode int

const (
	stickUnknown stickCode = iota
	stickUp
	stickRight
	stickDown
	stickLeft
)
