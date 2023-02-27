//go:build example

package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	input "github.com/quasilyte/ebitengine-input"
)

// Note: this example won't work well in browsers (wasm builds)
// due to the fact that browsers don't "connect" the gamepads until
// the user presses any button.
// See _examples/gamepad_in_browser for an example that works around it.

// Define our list of actions as enum-like constants.
//
// Actions usually have more than one key associated with them.
// A key could be a keyboard key, a gamepad key, a mouse button, etc.
//
// When you want to check whether the player is pressing the "fire" key,
// instead of checking the left mouse button directly, you check whether ActionFire is active.
const (
	ActionUnknown input.Action = iota
	ActionDebug
	ActionMoveLeft
	ActionMoveUp
	ActionMoveRight
	ActionMoveDown
	ActionExit
	ActionTeleport
)

func main() {
	ebiten.SetWindowSize(640, 480)

	if err := ebiten.RunGame(newExampleGame()); err != nil {
		log.Fatal(err)
	}
}

type exampleGame struct {
	started bool

	players []*player

	message string

	inputHandlers []*input.Handler
	inputSystem   input.System
}

func newExampleGame() *exampleGame {
	g := &exampleGame{}

	// The System.Init() should be called exactly once.
	g.inputSystem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})

	return g
}

func (g *exampleGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func (g *exampleGame) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, g.message)
	for _, p := range g.players {
		p.Draw(screen)
	}
}

func (g *exampleGame) Update() error {
	g.inputSystem.Update()

	if !g.started {
		g.Init()
		g.started = true
	}

	// Treat the first input handler as the main one.
	// Only the first player can exit the game.
	if g.inputHandlers[0].ActionIsJustPressed(ActionExit) {
		os.Exit(0)
	}
	if g.inputHandlers[0].ActionIsJustPressed(ActionDebug) {
		fmt.Println("debug action is pressed")
	}

	for _, p := range g.players {
		p.Update()
	}

	return nil
}

func (g *exampleGame) Init() {
	// We're hardcoding the keymap here, but it could be read from the config file.
	keymap := input.Keymap{
		// Every action can have a list of keys that can activate it.
		// KeyGamepadLStick<Direction> implements a D-pad like events for L/R sticks.
		ActionMoveLeft:  {input.KeyGamepadLeft, input.KeyGamepadLStickLeft, input.KeyLeft, input.KeyA},
		ActionMoveUp:    {input.KeyGamepadUp, input.KeyGamepadLStickUp, input.KeyUp, input.KeyW},
		ActionMoveRight: {input.KeyGamepadRight, input.KeyGamepadLStickRight, input.KeyRight, input.KeyD},
		ActionMoveDown:  {input.KeyGamepadDown, input.KeyGamepadLStickDown, input.KeyDown, input.KeyS},
		ActionExit: {
			input.KeyGamepadStart,
			input.KeyEscape,
			input.KeyWithModifier(input.KeyC, input.ModControl),
		},
		ActionDebug: {input.KeyControlLeft, input.KeyGamepadLStick, input.KeyGamepadRStick},
	}

	// Player 1 will have a teleport ability activated by a mouse click or a touch screen tap.
	keymap0 := keymap.Clone()
	keymap0[ActionTeleport] = []input.Key{input.KeyMouseLeft, input.KeyTouchTap}

	// Prepare the input handlers for all possible player slots.
	numGamepads := 0
	g.inputHandlers = make([]*input.Handler, 4)
	for i := range g.inputHandlers {
		m := keymap
		if i == 0 {
			m = keymap0
		}
		h := g.inputSystem.NewHandler(uint8(i), m)
		if h.GamepadConnected() {
			numGamepads++
		}
		g.inputHandlers[i] = h
	}

	// There can be only one player with keyboard.
	// There can be up to 4 players with gamepads.
	numPlayers := 1
	inputDevice := input.KeyboardDevice
	if numGamepads != 0 {
		inputDevice = input.GamepadDevice
		numPlayers = numGamepads
	}

	// Depending on the actual number of players, create
	// player objects and give them associated input handlers.
	g.players = make([]*player, numPlayers)
	pos := image.Point{X: 256, Y: 128}
	for i := range g.players {
		g.players[i] = &player{
			input: g.inputHandlers[i],
			pos:   pos,
			label: fmt.Sprintf("[player%d]", i+1),
		}
		pos.Y += 64
	}

	// For the real-world games you will want to map these action key names to
	// something more human-readable (you may also want to translate them).
	// For simplicity, we'll use them here as is.
	messageLines := []string{
		"preferred input: " + inputDevice.String(),
		"move left: " + strings.Join(g.inputHandlers[0].ActionKeyNames(ActionMoveLeft, inputDevice), " or "),
		"move right: " + strings.Join(g.inputHandlers[0].ActionKeyNames(ActionMoveRight, inputDevice), " or "),
	}
	g.message = strings.Join(messageLines, "\n")
}

type player struct {
	label string
	input *input.Handler
	pos   image.Point
}

func (p *player) Update() {
	if p.input.ActionIsPressed(ActionMoveLeft) {
		p.pos.X -= 4
	}
	if p.input.ActionIsPressed(ActionMoveUp) {
		p.pos.Y -= 4
	}
	if p.input.ActionIsPressed(ActionMoveRight) {
		p.pos.X += 4
	}
	if p.input.ActionIsPressed(ActionMoveDown) {
		p.pos.Y += 4
	}
	if info, ok := p.input.JustPressedActionInfo(ActionTeleport); ok {
		p.pos.X = int(info.Pos.X)
		p.pos.Y = int(info.Pos.Y)
	}
}

func (p *player) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, p.label, p.pos.X, p.pos.Y)
}
