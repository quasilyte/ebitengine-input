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

	// This is a scratch slice for ebiten.AppendPressedKeys operation.
	keySlice        []ebiten.Key
	gamepadKeySlice []ebiten.GamepadButton

	pendingEvents       []simulatedEvent
	prevSimulatedEvents []simulatedEvent
	simulatedEvents     []simulatedEvent
	hasSimulatedActions bool

	touchEnabled     bool
	touchHasTap      bool
	touchHasLongTap  bool
	touchJustHadDrag bool
	touchHasDrag     bool
	touchDragging    bool
	touchIDs         []ebiten.TouchID // This is a scratch slice, we don't support multi-touches yet
	touchActiveID    ebiten.TouchID
	touchTapPos      Vec
	touchDragPos     Vec
	touchStartPos    Vec
	touchTime        float64

	mouseEnabled          bool
	mouseHasDrag          bool // For "drag" event
	mouseDragging         bool // For "drag" event
	mouseJustHadDrag      bool // For "drag" event
	mouseJustReleasedDrag bool // For "drag" event
	mousePressed          bool // For "drag" event
	mouseStartPos         Vec  // For "drag" event
	mouseDragPos          Vec  // For "drag" event
	cursorPos             Vec
	wheel                 Vec
}

// SystemConfig configures the input system.
// This configuration can't be changed once created.
type SystemConfig struct {
	// DevicesEnabled selects the input devices that should be handled.
	// For the most cases, AnyDevice value is a good option.
	DevicesEnabled DeviceKind
}

func (sys *System) Init(config SystemConfig) {
	sys.keySlice = make([]ebiten.Key, 0, 4)
	sys.gamepadKeySlice = make([]ebiten.GamepadButton, 0, 2)

	sys.touchEnabled = config.DevicesEnabled&TouchDevice != 0
	sys.mouseEnabled = config.DevicesEnabled&MouseDevice != 0

	sys.gamepadIDs = make([]ebiten.GamepadID, 0, 8)
	sys.gamepadInfo = make([]gamepadInfo, 8)

	if sys.touchEnabled {
		sys.touchIDs = make([]ebiten.TouchID, 0, 8)
		sys.touchActiveID = -1
	}
}

// UpdateWithDelta is like Update(), but it allows you to specify the time delta.
func (sys *System) UpdateWithDelta(delta float64) {
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
				switch {
				case ebiten.IsStandardGamepadLayoutAvailable(id):
					info.model = gamepadStandard
				case isFirefox():
					info.model = guessFirefoxGamepadModel(int(id))
				default:
					info.model = guessGamepadModel(modelName)
				}
			}
			sys.updateGamepadInfo(id, info)
		}
	}

	if sys.touchEnabled {
		sys.touchHasTap = false
		sys.touchHasLongTap = false
		sys.touchHasDrag = false
		sys.touchJustHadDrag = false
		// Track the touch gesture release.
		// If it was a tap, set a flag.
		if sys.touchActiveID != -1 && inpututil.IsTouchJustReleased(sys.touchActiveID) {
			if !sys.touchDragging {
				if sys.touchTime >= 0.5 {
					sys.touchHasLongTap = true
				} else {
					sys.touchHasTap = true
				}
				sys.touchTapPos = sys.touchStartPos
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
				sys.touchTime += delta
				if vecDistance(sys.touchStartPos, currentPos) > 5 {
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
				sys.touchTime = 0
				break
			}
		}
	}

	if sys.mouseEnabled {
		x, y := ebiten.CursorPosition()
		sys.cursorPos = Vec{X: float64(x), Y: float64(y)}

		// We copy a lot from the touch-style drag gesture.
		// This is not mandatory as getting a cursor pos is much easier on PC.
		// But I do value the consistency and easier cross-platform coding,
		// so let's try to make them behave as close to each other as feasible.
		sys.mouseHasDrag = false
		sys.mouseJustHadDrag = false
		sys.mouseJustReleasedDrag = false
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			if sys.mouseDragging {
				sys.mouseJustReleasedDrag = true
			}
			sys.mouseDragging = false
			sys.mousePressed = false
		}
		if sys.mousePressed {
			if sys.mouseDragging {
				sys.mouseHasDrag = true
				sys.mouseDragPos = sys.cursorPos
			} else {
				// Mouse pointer is more precise than a finger gesture,
				// therefore we can have a lower threshold here.
				if vecDistance(sys.mouseStartPos, sys.cursorPos) > 3 {
					sys.mouseDragging = true
					sys.mouseJustHadDrag = true
					sys.mouseHasDrag = true
					sys.mouseDragPos = sys.cursorPos
				}
			}
		}
		if !sys.mousePressed && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			sys.mouseStartPos = sys.cursorPos
			sys.mousePressed = true
		}
	}

	if sys.mouseEnabled || sys.touchEnabled {
		x, y := ebiten.Wheel()
		sys.wheel = Vec{X: x, Y: y}
	}
}

// Update reads the input state and updates the information
// available to all input handlers.
// Generally, you call this method from your ebiten.Game.Update() method.
//
// Since ebitengine uses a fixed timestep architecture,
// a time delta of 1.0/60.0 is implied.
// If you need a control over that, use UpdateWithDelta() instead.
//
// The time delta mostly needed for things like press gesture
// detection: we need to calculate when a tap becomes a [long] press.
func (sys *System) Update() {
	sys.UpdateWithDelta(1.0 / 60.0)
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
			v := ebiten.GamepadAxisValue(id, axis)
			info.axisValues[axis] = v
		}
	}
}

// NewHandler creates a handler associated with player/device ID.
// IDs should start with 0 with a step of 1.
// So, NewHandler(0, ...) then NewHandler(1, ...).
//
// If you want to configure the handler further, use Handler fields/methods
// to do that. For example, see Handler.GamepadDeadzone.
func (sys *System) NewHandler(playerID uint8, keymap Keymap) *Handler {
	return &Handler{
		id:     playerID,
		keymap: keymap,
		sys:    sys,

		// My gamepads may have false positive activations with a
		// value lower than 0.03; we're using 0.055 here just to be safe.
		// Various sources indicate that a value of ~0.05 is optimal for a default.
		GamepadDeadzone: 0.055,
	}
}
