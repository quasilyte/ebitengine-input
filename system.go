package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// System is the main component of the input library.
//
// You usually need only one input system object.
//
// Store System object (by value) inside your game context/state object like this:
//
//	struct GameState {
//	    InputSystem input.System
//	}
//
// When ebitengine game is executed, call gameState.InputSystem.Init() once.
//
// On every ebitengine Update() call, use gameState.InputSystem.Update().
//
// The system is usually not used directly after the input handlers are created.
// Use input handlers to handle the user input.
type System struct {
	gamepadIDs  []ebiten.GamepadID
	gamepadInfo []gamepadInfo

	pendingEvents       []simulatedEvent
	prevSimulatedEvents []simulatedEvent
	simulatedEvents     []simulatedEvent
	hasSimulatedActions bool

	touchEnabled     bool
	touchHasTap      bool
	touchJustHadDrag bool
	touchHasDrag     bool
	touchDragging    bool
	touchIDs         []ebiten.TouchID // This is a scratch slice, we don't support multi-touches yet
	touchActiveID    ebiten.TouchID
	touchTapPos      Vec
	touchDragPos     Vec
	touchStartPos    Vec

	mouseEnabled bool
	cursorPos    Vec
	wheel        Vec
}

// SystemConfig configures the input system.
// This configuration can't be changed once created.
type SystemConfig struct {
	// DevicesEnabled selects the input devices that should be handled.
	// For the most cases, AnyDevice value is a good option.
	DevicesEnabled DeviceKind
}

func (sys *System) Init(config SystemConfig) {
	sys.touchEnabled = config.DevicesEnabled&TouchDevice != 0
	sys.mouseEnabled = config.DevicesEnabled&MouseDevice != 0

	sys.gamepadIDs = make([]ebiten.GamepadID, 0, 8)
	sys.gamepadInfo = make([]gamepadInfo, 8)

	if sys.touchEnabled {
		sys.touchIDs = make([]ebiten.TouchID, 0, 8)
		sys.touchActiveID = -1
	}
}

func (sys *System) Update() {
	// Rotate the events slices.
	// Pending events become simulated in this frame.
	// Re-use the other slice capacity to push new events.
	//	prev simulated <- simulated
	//	pending <- prev simulated
	//	simulated <- pending
	sys.prevSimulatedEvents, sys.pendingEvents, sys.simulatedEvents =
		sys.simulatedEvents, sys.prevSimulatedEvents, sys.pendingEvents
	sys.pendingEvents = sys.pendingEvents[:0]
	sys.hasSimulatedActions = false
	for i := range sys.simulatedEvents {
		if sys.simulatedEvents[i].keyKind == keySimulated {
			sys.hasSimulatedActions = true
			break
		}
	}

	sys.gamepadIDs = ebiten.AppendGamepadIDs(sys.gamepadIDs[:0])
	if len(sys.gamepadIDs) != 0 {
		for i, id := range sys.gamepadIDs {
			info := &sys.gamepadInfo[i]
			info.axisCount = ebiten.GamepadAxisCount(id)
			modelName := ebiten.GamepadName(id)
			if info.modelName != modelName {
				info.modelName = modelName
				if ebiten.IsStandardGamepadLayoutAvailable(id) {
					info.model = gamepadStandard
				} else if isFirefox() {
					info.model = guessFirefoxGamepadModel(int(id))
				} else {
					info.model = guessGamepadModel(modelName)
				}
			}
			sys.updateGamepadInfo(id, info)
		}
	}

	if sys.touchEnabled {
		sys.touchHasTap = false
		sys.touchHasDrag = false
		sys.touchJustHadDrag = false
		// Track the touch gesture release.
		// If it was a tap, set a flag.
		if sys.touchActiveID != -1 && inpututil.IsTouchJustReleased(sys.touchActiveID) {
			if !sys.touchDragging {
				x, y := ebiten.TouchPosition(sys.touchActiveID)
				releasePos := Vec{X: float64(x), Y: float64(y)}
				sys.touchHasTap = true
				sys.touchTapPos = releasePos
			}
			sys.touchActiveID = -1
			sys.touchDragging = false
		}
		// Check if this gesture entered a drag mode.
		// Drag mode gestures will not trigger a tap when released.
		// Drag events emit a pos delta relative to a start pos every frame.
		if sys.touchActiveID != -1 {
			x, y := ebiten.TouchPosition(sys.touchActiveID)
			currentPos := Vec{X: float64(x), Y: float64(y)}
			if sys.touchDragging {
				sys.touchHasDrag = true
				sys.touchDragPos = currentPos
			} else {
				if vecDistance(sys.touchStartPos, currentPos) > 3 {
					sys.touchDragging = true
					sys.touchJustHadDrag = true
					sys.touchHasDrag = true
					sys.touchDragPos = currentPos
				}
			}
		}
		// Check if a new touch gesture is started.
		if sys.touchActiveID == -1 {
			sys.touchIDs = inpututil.AppendJustPressedTouchIDs(sys.touchIDs[:0])
			for _, id := range sys.touchIDs {
				x, y := ebiten.TouchPosition(id)
				sys.touchStartPos = Vec{X: float64(x), Y: float64(y)}
				sys.touchActiveID = id
				break
			}
		}
	}

	if sys.mouseEnabled {
		x, y := ebiten.CursorPosition()
		sys.cursorPos = Vec{X: float64(x), Y: float64(y)}
	}

	if sys.mouseEnabled || sys.touchEnabled {
		x, y := ebiten.Wheel()
		sys.wheel = Vec{X: x, Y: y}
	}
}

func (sys *System) updateGamepadInfo(id ebiten.GamepadID, info *gamepadInfo) {
	switch info.model {
	case gamepadStandard:
		copy(info.prevAxisValues[:], info.axisValues[:])
		for axis := ebiten.StandardGamepadAxisLeftStickHorizontal; axis <= ebiten.StandardGamepadAxisMax; axis++ {
			v := ebiten.StandardGamepadAxisValue(id, axis)
			info.axisValues[int(axis)] = v
		}
	case gamepadFirefoxXinput:
		copy(info.prevAxisValues[:], info.axisValues[:])
		for axis := 0; axis < info.axisCount; axis++ {
			axisKey := axis
			switch ebiten.StandardGamepadAxis(axis) {
			case ebiten.StandardGamepadAxisLeftStickHorizontal:
				axisKey = 0
			case ebiten.StandardGamepadAxisLeftStickVertical:
				axisKey = 1
			case ebiten.StandardGamepadAxisRightStickHorizontal:
				axisKey = 3
			case ebiten.StandardGamepadAxisRightStickVertical:
				axisKey = 4
			}
			v := ebiten.GamepadAxisValue(id, axisKey)
			info.axisValues[axis] = v
		}
	}
}

// NewHandler creates a handler associated with player/device ID.
// IDs should start with 0 with a step of 1.
// So, NewHandler(0, ...) then NewHandler(1, ...).
func (sys *System) NewHandler(playerID uint8, keymap Keymap) *Handler {
	return &Handler{
		id:     playerID,
		keymap: keymap,
		sys:    sys,
	}
}
