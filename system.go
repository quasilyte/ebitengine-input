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
//    struct GameState {
//        InputSystem input.System
//    }
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

	touchEnabled bool
	touchIDs     []ebiten.TouchID
	touchTapID   ebiten.TouchID
	touchHasTap  bool
	touchTapPos  Point

	mouseEnabled bool
	cursorPos    Point
}

// SystemConfig configures the input system.
// This configuration can't be changed once created.
type SystemConfig struct {
	// DevicesEnabled selects the input devices that should be handled.
	// For the most cases, AnyInput value is a good option.
	DevicesEnabled InputDeviceKind
}

func (sys *System) Init(config SystemConfig) {
	sys.touchEnabled = config.DevicesEnabled&TouchInput != 0
	sys.mouseEnabled = config.DevicesEnabled&MouseInput != 0

	sys.gamepadIDs = make([]ebiten.GamepadID, 0, 8)
	sys.gamepadInfo = make([]gamepadInfo, 8)

	if sys.touchEnabled {
		sys.touchIDs = make([]ebiten.TouchID, 0, 8)
	}
}

func (sys *System) Update() {
	sys.gamepadIDs = ebiten.AppendGamepadIDs(sys.gamepadIDs[:0])
	if len(sys.gamepadIDs) != 0 {
		for i, id := range sys.gamepadIDs {
			info := &sys.gamepadInfo[i]
			modelName := ebiten.GamepadName(id)
			if info.modelName != modelName {
				info.modelName = modelName
				if ebiten.IsStandardGamepadLayoutAvailable(id) {
					info.model = gamepadStandard
				} else {
					info.model = guessGamepadModel(modelName)
				}
			}
		}
	}

	if sys.touchEnabled {
		sys.touchHasTap = false
		for _, id := range sys.touchIDs {
			if id == sys.touchTapID && inpututil.IsTouchJustReleased(id) {
				sys.touchHasTap = true
				break
			}
		}
		if !sys.touchHasTap {
			sys.touchIDs = inpututil.AppendJustPressedTouchIDs(sys.touchIDs)
			for _, id := range sys.touchIDs {
				x, y := ebiten.TouchPosition(id)
				sys.touchTapPos = Point{X: float64(x), Y: float64(y)}
				sys.touchTapID = id
				break
			}
		}
	}

	if sys.mouseEnabled {
		x, y := ebiten.CursorPosition()
		sys.cursorPos = Point{X: float64(x), Y: float64(y)}
	}
}

// NewHandler creates a handler associated with player/device ID.
// IDs should start with 0 with a step of 1.
// So, NewHandler(0, ...) then NewHandler(1, ...).
func (sys *System) NewHandler(playerID int, keymap Keymap) *Handler {
	return &Handler{
		id:     playerID,
		keymap: keymap,
		sys:    sys,
	}
}

func (sys *System) NewMultiHandler() *MultiHandler {
	return &MultiHandler{}
}
